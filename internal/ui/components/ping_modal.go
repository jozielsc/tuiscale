package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"tuiscale/internal/ui/theme"
)

// RenderPingModal desenha o modal com o resultado do comando `tailscale ping`.
func RenderPingModal(targetHost, targetIP string, lines []string, isRunning bool, width, height int) string {
	modalWidth := 70
	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	title := fmt.Sprintf("PING: %s (%s)", targetHost, targetIP)
	header := theme.ModalTitle.Render(title)

	var bodyLines []string
	if len(lines) == 0 && isRunning {
		bodyLines = append(bodyLines, "Disparando pacotes ping Tailscale... Aguarde...")
	} else {
		for _, line := range lines {
			if strings.Contains(line, "via direct") || strings.Contains(line, "via 1") || strings.Contains(line, "via [") {
				bodyLines = append(bodyLines, theme.BadgeDirect.Render(line))
			} else if strings.Contains(line, "via DERP") {
				bodyLines = append(bodyLines, theme.BadgeRelay.Render(line))
			} else if strings.Contains(line, "timed out") || strings.Contains(line, "error") {
				bodyLines = append(bodyLines, theme.ToastError.Render(line))
			} else {
				bodyLines = append(bodyLines, lipgloss.NewStyle().Foreground(theme.ColorText).Render(line))
			}
		}
	}

	hint := lipgloss.NewStyle().Foreground(theme.ColorTextDim).Render("\nPressione [Esc] ou [Enter] para fechar")
	content := fmt.Sprintf("%s\n\n%s\n%s", header, strings.Join(bodyLines, "\n"), hint)

	renderedBox := theme.ModalBox.Width(modalWidth).Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, renderedBox)
}
