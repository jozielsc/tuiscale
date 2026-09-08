package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/tailscale"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// RenderPeerDetail desenha o painel de inspeção minucioso de um peer específico.
func RenderPeerDetail(p *tailscale.PeerStatus, width, height int) string {
	if p == nil {
		return theme.Panel.Width(width).Height(height).Render("Selecione um peer para ver os detalhes.")
	}

	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorAccent)
	valStyle := lipgloss.NewStyle().Foreground(theme.ColorText)
	dimValStyle := lipgloss.NewStyle().Foreground(theme.ColorTextMuted)

	// Espaço interno disponível para valores (descontando label ~22, bordas e padding)
	labelWidth := 22
	maxVal := width - labelWidth - 4
	if maxVal < 8 {
		maxVal = 8
	}

	// formatRow para valores de TEXTO PURO (sem ANSI). Aplica truncamento seguro.
	formatRow := func(label, val string) string {
		val = truncate(val, maxVal)
		return fmt.Sprintf("%-22s %s", labelStyle.Render(label+":"), valStyle.Render(val))
	}

	// formatRowStyled para valores JÁ ESTILIZADOS (contêm ANSI). Não trunca para
	// não corromper as sequências de escape — esses valores são curtos por natureza.
	formatRowStyled := func(label, styledVal string) string {
		return fmt.Sprintf("%-22s %s", labelStyle.Render(label+":"), styledVal)
	}

	var statusTag string
	if p.Active {
		statusTag = theme.BadgeDirect.Render("🟢 Ativo (Sessão direta)")
	} else if p.Online {
		statusTag = theme.BadgeRelay.Render("🔵 Online (Em espera ou via relay)")
	} else {
		statusTag = theme.BadgeOffline.Render("⚪ Offline")
	}

	ipsStr := strings.Join(p.TailscaleIPs, ", ")
	if ipsStr == "" {
		ipsStr = "-"
	}

	routeStr := "Nenhuma"
	if p.CurAddr != "" {
		routeStr = fmt.Sprintf("Direta UDP (%s)", p.CurAddr)
	} else if p.Relay != "" {
		city, country := tailscale.GetDERPLocationByCode(p.Relay)
		routeStr = fmt.Sprintf("Relay DERP [%s] (%s, %s)", p.Relay, city, country)
	}

	exitNodeStr := "Não"
	if p.ExitNode {
		exitNodeStr = theme.BadgeExitNode.Render("ATIVO (Roteando tráfego)")
	} else if p.ExitNodeOption {
		exitNodeStr = "Disponível (Pode ser ativado)"
	}

	lines := []string{
		theme.TitleStyle.Render("DETALHES DO DISPOSITIVO"),
		"",
		formatRow("Nome do Host", p.HostName),
		formatRow("DNS MagicDNS", strings.TrimSuffix(p.DNSName, ".")),
		formatRowStyled("Sistema Operacional", theme.FormatOSIcon(strings.ToLower(p.OS))),
		formatRowStyled("Status de Conexão", statusTag),
		formatRow("Endereços Tailscale", ipsStr),
		formatRow("Usuário Proprietário", p.UserLogin),
		formatRow("Rota de Rede Atual", routeStr),
		formatRowStyled("Nó de Saída (Exit Node)", exitNodeStr),
		"",
		theme.SubTitleStyle.Render("Tráfego de Dados:"),
		formatRow("Total Recebido (Rx)", tailscale.FormatBytes(p.RxBytes)),
		formatRow("Total Enviado (Tx)", tailscale.FormatBytes(p.TxBytes)),
		formatRowStyled("Taxa Instantânea", fmt.Sprintf("↓ %s   ↑ %s",
			tailscale.FormatSpeed(p.SpeedRx),
			tailscale.FormatSpeed(p.SpeedTx),
		)),
		"",
		theme.SubTitleStyle.Render("Segurança e Chaves:"),
		formatRow("ID do Nó", fmt.Sprintf("%d", p.NodeID)),
		formatRow("Expiração da Chave", p.KeyExpiry),
		formatRow("Último Handshake", p.LastHandshake),
		formatRow("Visto pela última vez", p.LastSeen),
	}

	// Adicionar endpoints físicos descobertos se existirem
	if len(p.Addrs) > 0 {
		lines = append(lines, "", theme.SubTitleStyle.Render("Endpoints Físicos (STUN/NAT):"))
		maxAddr := width - 6
		if maxAddr < 8 {
			maxAddr = 8
		}
		for _, addr := range p.Addrs {
			lines = append(lines, dimValStyle.Render("  • "+truncate(addr, maxAddr)))
		}
	}

	content := strings.Join(lines, "\n")
	return theme.Panel.Width(width).Height(height).Render(content)
}
