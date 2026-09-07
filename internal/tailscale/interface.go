package tailscale

import "context"

// TailscaleClient define a interface para operações do cliente Tailscale.
// Esta interface permite criar mocks para testes unitários.
type TailscaleClient interface {
	// GetStatus retorna o estado atual do Tailscale.
	GetStatus(ctx context.Context) (*Status, error)

	// RunNetcheck executa diagnóstico de rede e retorna o relatório.
	RunNetcheck(ctx context.Context) (*NetcheckReport, error)

	// PingPeer executa ping contra um peer e retorna as linhas de resposta.
	PingPeer(ctx context.Context, target string, count int) ([]string, error)

	// Connect conecta o Tailscale à rede.
	Connect(ctx context.Context) error

	// Disconnect desconecta o Tailscale da rede.
	Disconnect(ctx context.Context) error

	// SetExitNode configura ou desativa o exit-node.
	SetExitNode(ctx context.Context, exitNode string) error

	// GetSpeedTracker retorna o rastreador de velocidade de tráfego.
	GetSpeedTracker() *SpeedTracker
}

// Garantir que Client implementa TailscaleClient
var _ TailscaleClient = (*Client)(nil)
