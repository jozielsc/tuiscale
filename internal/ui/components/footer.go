package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// RenderFooter desenha a barra de rodapé com atalhos de teclado e mensagens de status (toast).
// Em telas estreitas os atalhos são reduzidos para os essenciais; o footer nunca some.
func RenderFooter(toastMsg string, isToastError bool, filterMode bool, filterText string, width int) string {
	if filterMode {
		prompt := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("Buscar: ")
		cursor := lipgloss.NewStyle().Foreground(theme.ColorAccent).Render("█")
		text := lipgloss.NewStyle().Foreground(theme.ColorText).Render(filterText)
		hint := ""
		if width >= 70 {
			hint = lipgloss.NewStyle().Foreground(theme.ColorTextDim).Render("  [Enter]/[Esc] para finalizar")
		}
		return theme.FooterBar.Width(width).Render(prompt + text + cursor + hint)
	}

	renderKey := func(k, desc string) string {
		return fmt.Sprintf("%s %s", theme.KeyShortcut.Render("["+k+"]"), theme.KeyDesc.Render(desc))
	}

	// Conjunto completo de atalhos (telas >= 100 colunas)
	fullShortcuts := strings.Join([]string{
		renderKey("c", "Conectar"),
		renderKey("d", "Desconectar"),
		renderKey("p", "Ping"),
		renderKey("e", "Exit Node"),
		renderKey("y", "Copiar IP"),
		renderKey("/", "Buscar"),
		renderKey("n", "Netcheck"),
		renderKey("?", "Ajuda"),
		renderKey("q", "Sair"),
	}, "  ")

	// Conjunto reduzido (telas < 100 colunas)
	shortShortcuts := strings.Join([]string{
		renderKey("c", "Up"),
		renderKey("d", "Down"),
		renderKey("p", "Ping"),
		renderKey("/", "Buscar"),
		renderKey("?", "Ajuda"),
		renderKey("q", "Sair"),
	}, "  ")

	// Ultra-compacto: apenas os mais críticos (telas muito pequenas)
	minShortcuts := strings.Join([]string{
		renderKey("?", "Ajuda"),
		renderKey("q", "Sair"),
	}, "  ")

	var shortcutsStr string
	switch {
	case width >= 100:
		shortcutsStr = fullShortcuts
	case width >= 60:
		shortcutsStr = shortShortcuts
	default:
		shortcutsStr = minShortcuts
	}

	var toastStr string
	if toastMsg != "" {
		if isToastError {
			toastStr = theme.ToastError.Render("✗ " + toastMsg)
		} else {
			toastStr = theme.ToastSuccess.Render("✓ " + toastMsg)
		}
	}

	if toastStr != "" {
		gap := width - lipgloss.Width(shortcutsStr) - lipgloss.Width(toastStr) - 3
		if gap >= 1 {
			content := shortcutsStr + strings.Repeat(" ", gap) + toastStr
			return theme.FooterBar.Width(width).Render(content)
		}
		// Sem espaço para os dois: exibe só o toast
		return theme.FooterBar.Width(width).Render(toastStr)
	}

	return theme.FooterBar.Width(width).Render(shortcutsStr)
}
