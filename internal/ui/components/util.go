package components

import "github.com/charmbracelet/lipgloss"

// truncate limita uma string a maxLen colunas de terminal, adicionando "…" se truncada.
// Usa lipgloss.Width para contar corretamente caracteres multibyte e emojis.
func truncate(s string, maxLen int) string {
	if maxLen < 1 {
		return ""
	}
	if lipgloss.Width(s) <= maxLen {
		return s
	}
	// Percorre rune a rune para não cortar no meio de um multibyte
	runes := []rune(s)
	w := 0
	cut := 0
	for i, r := range runes {
		rw := lipgloss.Width(string(r))
		if w+rw > maxLen-1 { // -1 reserva espaço para "…"
			cut = i
			break
		}
		w += rw
		cut = i + 1
	}
	return string(runes[:cut]) + "…"
}
