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
type Client struct {
	mu           sync.RWMutex
	tailscaleBin string
	socketPath   string
	speedTracker *SpeedTracker
	lastStatus   *Status
	lastNetcheck *NetcheckReport
}

// NewClient cria uma nova instância de Client.
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
	return lines, err
}

// Connect aciona `tailscale up`.
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

// Disconnect aciona `tailscale down`.
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
