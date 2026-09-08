package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/tailscale"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// RenderHeader desenha o banner superior com informações da máquina local, status e velocidade.
func RenderHeader(status *tailscale.Status, speedTracker *tailscale.SpeedTracker, width int) string {
	if status == nil {
		return theme.HeaderBox.Width(width - 2).Render("Carregando status do Tailscale...")
	}

	// 1. Badge de Estado
	var statusBadge string
	switch strings.ToLower(status.BackendState) {
	case "running":
		statusBadge = theme.BadgeRunning.Render("● CONECTADO")
	case "stopped":
		statusBadge = theme.BadgeStopped.Render("■ DESCONECTADO")
	case "needslogin":
		statusBadge = theme.BadgeStarting.Render("▲ LOGIN NECESSÁRIO")
	case "starting":
		statusBadge = theme.BadgeStarting.Render("◌ INICIANDO...")
	default:
		statusBadge = theme.BadgeStopped.Render(strings.ToUpper(status.BackendState))
	}

	// 2. Dados do Nó Local (Self)
	var localHost, localIP, localDNS, account string
	if status.Self != nil {
		localHost = status.Self.HostName
		if len(status.Self.TailscaleIPs) > 0 {
			localIP = status.Self.TailscaleIPs[0]
		}
		localDNS = strings.TrimSuffix(status.Self.DNSName, ".")
		account = status.Self.UserLogin
	}
	if account == "" && status.CurrentTailnet != nil {
		account = status.CurrentTailnet.Name
	}

	// 3. Velocidade em Tempo Real e Tráfego Acumulado
	var rxSpeedStr, txSpeedStr, totalRxStr, totalTxStr string
	if speedTracker != nil {
		rxSpeedStr = tailscale.FormatSpeed(speedTracker.CurrentRx)
		txSpeedStr = tailscale.FormatSpeed(speedTracker.CurrentTx)
		totalRxStr = tailscale.FormatBytes(speedTracker.CumulativeRx)
		totalTxStr = tailscale.FormatBytes(speedTracker.CumulativeTx)
	} else {
		rxSpeedStr = "0 B/s"
		txSpeedStr = "0 B/s"
		totalRxStr = "0 B"
		totalTxStr = "0 B"
	}

	speedBlock := fmt.Sprintf("Velocidade: %s  %s  |  Total: ↓ %s  ↑ %s",
		theme.SpeedRxStyle.Render("↓ "+rxSpeedStr),
		theme.SpeedTxStyle.Render("↑ "+txSpeedStr),
		totalRxStr,
		totalTxStr,
	)

	// Linha 1: Título do App, Badge de Conexão, Conta
	leftLine1 := fmt.Sprintf("%s  %s",
		theme.TitleStyle.Render("TUIScale"),
		statusBadge,
	)
	rightLine1 := ""
	if account != "" {
		rightLine1 = fmt.Sprintf("%s: %s",
			theme.SubTitleStyle.Render("Conta"),
			lipgloss.NewStyle().Foreground(theme.ColorAccent).Render(account),
		)
	}

	// Linha 2: Host, IP local, DNS e Velocidade
	leftLine2 := fmt.Sprintf("%s: %s  |  %s: %s",
		theme.SubTitleStyle.Render("Máquina"),
		lipgloss.NewStyle().Bold(true).Foreground(theme.ColorText).Render(localHost),
		theme.SubTitleStyle.Render("IP"),
		lipgloss.NewStyle().Foreground(theme.ColorHighlight).Render(localIP),
	)
	if localDNS != "" && width > 110 {
		leftLine2 += fmt.Sprintf("  |  %s: %s", theme.SubTitleStyle.Render("MagicDNS"), theme.SubTitleStyle.Render(localDNS))
	}

	// Ajuste de largura e alinhamento
	innerWidth := width - 4
	if innerWidth < 40 {
		innerWidth = 40
	}

	gap1 := innerWidth - lipgloss.Width(leftLine1) - lipgloss.Width(rightLine1)
	if gap1 < 1 {
		gap1 = 1
	}
	line1 := leftLine1 + strings.Repeat(" ", gap1) + rightLine1

	gap2 := innerWidth - lipgloss.Width(leftLine2) - lipgloss.Width(speedBlock)
	if gap2 < 1 {
		gap2 = 1
	}
	line2 := leftLine2 + strings.Repeat(" ", gap2) + speedBlock

	content := line1 + "\n" + line2
	return theme.HeaderBox.Width(width - 2).Render(content)
}
