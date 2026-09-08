package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// RenderFooter desenha a barra de rodapé com atalhos de teclado e mensagens de status (toast).
func RenderFooter(toastMsg string, isToastError bool, filterMode bool, filterText string, width int) string {
	if filterMode {
		prompt := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("Buscar: ")
		cursor := lipgloss.NewStyle().Blink(true).Foreground(theme.ColorAccent).Render("█")
		text := lipgloss.NewStyle().Foreground(theme.ColorText).Render(filterText)
		hint := lipgloss.NewStyle().Foreground(theme.ColorTextDim).Render("  (Pressione [Enter] ou [Esc] para finalizar)")
		return theme.FooterBar.Width(width).Render(prompt + text + cursor + hint)
	}

	// Atalhos rápidos
	renderKey := func(k, desc string) string {
		return fmt.Sprintf("%s %s", theme.KeyShortcut.Render("["+k+"]"), theme.KeyDesc.Render(desc))
	}

	shortcuts := []string{
		renderKey("c", "Conectar"),
		renderKey("d", "Desconectar"),
		renderKey("p", "Ping"),
		renderKey("e", "Exit Node"),
		renderKey("y", "Copiar IP"),
		renderKey("/", "Buscar"),
		renderKey("n", "Netcheck"),
		renderKey("?", "Ajuda"),
		renderKey("q", "Sair"),
	}

	shortcutsStr := strings.Join(shortcuts, "  ")

	var toastStr string
	if toastMsg != "" {
		if isToastError {
			toastStr = theme.ToastError.Render("✗ " + toastMsg)
		} else {
			toastStr = theme.ToastSuccess.Render("✓ " + toastMsg)
		}
	}

	// Layout responsivo para rodapé
	availableSpace := width - lipgloss.Width(shortcutsStr) - 4
	if toastStr != "" && availableSpace > 10 {
		gap := width - lipgloss.Width(shortcutsStr) - lipgloss.Width(toastStr) - 3
		if gap < 1 {
			gap = 1
		}
		content := shortcutsStr + strings.Repeat(" ", gap) + toastStr
		return theme.FooterBar.Width(width).Render(content)
	}

	if toastStr != "" && width < 90 {
		return theme.FooterBar.Width(width).Render(toastStr)
	}

	return theme.FooterBar.Width(width).Render(shortcutsStr)
}
