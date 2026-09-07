package tailscale

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Client gerencia a comunicação com o Tailscale CLI e daemon.
// Ele fornece métodos para obter status, executar diagnósticos de rede,
// e controlar o estado da conexão Tailscale.
type Client struct {
	mu           sync.RWMutex
	tailscaleBin string
	socketPath   string
	speedTracker *SpeedTracker
	lastStatus   *Status
	lastNetcheck *NetcheckReport
}

// NewClient cria uma nova instância de Client.
// Se socketPath for vazio, usa o socket padrão do sistema (/var/run/tailscale/tailscaled.sock).
func NewClient(socketPath string) *Client {
	bin, err := exec.LookPath("tailscale")
	if err != nil {
		bin = "tailscale"
	}
	return &Client{
		tailscaleBin: bin,
		socketPath:   socketPath,
		speedTracker: NewSpeedTracker(),
	}
}

// buildCommand constrói a execução do comando tailscale adicionando flag de socket se fornecida.
func (c *Client) buildCommand(ctx context.Context, args ...string) *exec.Cmd {
	var finalArgs []string
	if c.socketPath != "" {
		finalArgs = append(finalArgs, "--socket="+c.socketPath)
	}
	finalArgs = append(finalArgs, args...)
	return exec.CommandContext(ctx, c.tailscaleBin, finalArgs...)
}

// GetStatus executa `tailscale status --json` e retorna o estado atual.
// Retorna erro se o daemon não estiver rodando ou se a resposta for inválida.
func (c *Client) GetStatus(ctx context.Context) (*Status, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	cmd := c.buildCommand(ctxTimeout, "status", "--json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return nil, fmt.Errorf("%s", errMsg)
		}
		return nil, err
	}

	raw := stdout.Bytes()
	// Localiza início do JSON em caso de logs impressos antes
	startIdx := bytes.IndexByte(raw, '{')
	if startIdx == -1 {
		return nil, errors.New("resposta do tailscale não contém JSON válido")
	}

	var status Status
	if err := json.Unmarshal(raw[startIdx:], &status); err != nil {
		return nil, fmt.Errorf("erro ao decodificar status: %w", err)
	}

	// Associar nomes de usuário e logins aos peers
	for _, peer := range status.Peer {
		if peer == nil {
			continue
		}
		if u, ok := status.User[strconv.FormatInt(peer.UserID, 10)]; ok {
			peer.UserLogin = u.LoginName
			if peer.UserLogin == "" {
				peer.UserLogin = u.DisplayName
			}
		}
	}
	if status.Self != nil {
		if u, ok := status.User[strconv.FormatInt(status.Self.UserID, 10)]; ok {
			status.Self.UserLogin = u.LoginName
		}
	}

	// Atualizar medição de velocidade
	c.speedTracker.Update(&status)

	c.mu.Lock()
	c.lastStatus = &status
	c.mu.Unlock()

	return &status, nil
}

// GetSpeedTracker retorna a instância de rastreamento de velocidade.
func (c *Client) GetSpeedTracker() *SpeedTracker {
	return c.speedTracker
}

// RunNetcheck executa `tailscale netcheck --format=json` e retorna o diagnóstico de rede.
// O netcheck verifica conectividade UDP, IPv4, IPv6, latência para DERP servers e tipo de NAT.
func (c *Client) RunNetcheck(ctx context.Context) (*NetcheckReport, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cmd := c.buildCommand(ctxTimeout, "netcheck", "--format=json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return nil, fmt.Errorf("%s", errMsg)
		}
		return nil, err
	}

	raw := stdout.Bytes()
	startIdx := bytes.IndexByte(raw, '{')
	lastIdx := bytes.LastIndexByte(raw, '}')
	if startIdx == -1 || lastIdx == -1 || lastIdx <= startIdx {
		return nil, errors.New("resposta de netcheck não contém JSON válido")
	}

	var report NetcheckReport
	if err := json.Unmarshal(raw[startIdx:lastIdx+1], &report); err != nil {
		return nil, fmt.Errorf("erro ao decodificar netcheck: %w", err)
	}

	c.mu.Lock()
	c.lastNetcheck = &report
	c.mu.Unlock()

	return &report, nil
}

// PingPeer executa `tailscale ping` contra um peer e retorna as linhas de resposta.
// O parâmetro count define o número de pings a enviar (padrão 4 se <= 0).
//
// Retorna as linhas de output mesmo em caso de erro, pois o tailscale ping
// pode retornar informações úteis mesmo quando o ping falha (ex: peer offline).
// O erro indica falha na execução do comando, não falha do ping em si.
func (c *Client) PingPeer(ctx context.Context, target string, count int) ([]string, error) {
	if count <= 0 {
		count = 4
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := c.buildCommand(ctxTimeout, "ping", fmt.Sprintf("--c=%d", count), "--until-direct=false", target)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	if output == "" && stderr.Len() > 0 {
		output = stderr.String()
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Filtrar linhas vazias
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	return result, err
}

// Connect aciona `tailscale up` para conectar à rede Tailscale.
// Requer que o usuário tenha permissão de operador (configurada via `tailscale set --operator=$USER`).
func (c *Client) Connect(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := c.buildCommand(ctxTimeout, "up", "--reset=false")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr == "" {
			errStr = strings.TrimSpace(stdout.String())
		}
		return fmt.Errorf("falha ao conectar: %s", errStr)
	}
	return nil
}

// Disconnect aciona `tailscale down` para desconectar da rede Tailscale.
// A conexão pode ser reestabelecida posteriormente com Connect().
func (c *Client) Disconnect(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()

	cmd := c.buildCommand(ctxTimeout, "down")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr == "" {
			errStr = strings.TrimSpace(stdout.String())
		}
		return fmt.Errorf("falha ao desconectar: %s", errStr)
	}
	return nil
}

// SetExitNode configura ou desativa o exit-node da máquina.
// Passar string vazia ("") desativa o exit-node atual.
// O exit-node deve ser um peer que oferece essa funcionalidade (ExitNodeOption=true).
func (c *Client) SetExitNode(ctx context.Context, exitNode string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()

	cmd := c.buildCommand(ctxTimeout, "set", "--exit-node="+exitNode)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr == "" {
			errStr = strings.TrimSpace(stdout.String())
		}
		return fmt.Errorf("falha ao configurar exit node: %s", errStr)
	}
	return nil
}

// ToggleShieldsUp alterna a proteção de conexões de entrada.
func (c *Client) ToggleShieldsUp(ctx context.Context, enable bool) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()

	val := "false"
	if enable {
		val = "true"
	}
	cmd := c.buildCommand(ctxTimeout, "set", "--shields-up="+val)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errStr := strings.TrimSpace(stderr.String())
		return fmt.Errorf("falha ao alterar shields-up: %s", errStr)
	}
	return nil
}

// IsDaemonNotRunningError verifica se o erro retornado indica que o daemon tailscaled está inativo ou inacessível.
func IsDaemonNotRunningError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	patterns := []string{
		"failed to connect to local tailscaled",
		"failed to connect to local tailscale daemon",
		"dial unix",
		"no such file or directory",
		"connection refused",
		"not running",
		"is tailscaled running",
		"cannot connect to tailscaled",
		"tailscale socket",
	}
	for _, p := range patterns {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}

