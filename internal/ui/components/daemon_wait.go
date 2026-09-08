package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// RenderDaemonWaitView desenha a tela amigável informando que o daemon tailscaled está inativo.
// O layout é totalmente responsivo: colapsa gradualmente à medida que width/height diminuem,
// garantindo compatibilidade com janelas pequenas como o scratchpad do Sway.
func RenderDaemonWaitView(width, height int, lastError string, attempts int) string {
	// Largura interna da caixa: usa o espaço disponível com mínimo defensivo
	boxWidth := width - 4
	if boxWidth > 74 {
		boxWidth = 74
	}
	if boxWidth < 20 {
		boxWidth = 20
	}

	// Tela muito pequena: exibe versão ultra-compacta de uma linha
	if width < 30 || height < 6 {
		msg := "tailscaled inativo  [q] Sair  [r] Retry"
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(theme.ColorWarning).Render(msg))
	}

	// ── 1. Título + Badge ──────────────────────────────────────────────────────
	var badge, descText string
	if attempts == 0 && lastError == "" {
		badge = theme.BadgeStarting.Render("⏳ CONECTANDO AO DAEMON")
		descText = "Iniciando comunicação com o daemon local do Tailscale..."
	} else {
		badge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(theme.ColorWarning).
			Padding(0, 1).
			Render("⚠️  SERVIÇO TAILSCALE INATIVO")
		descText = "Não foi possível se comunicar com o daemon local ('tailscaled').\nO serviço parece estar desligado ou ainda não foi inicializado."
	}

	title := lipgloss.JoinHorizontal(lipgloss.Center,
		theme.TitleStyle.Render("TUIScale"), "  ", badge,
	)

	// ── 2. Descrição ──────────────────────────────────────────────────────────
	desc := lipgloss.NewStyle().Foreground(theme.ColorText).Width(boxWidth - 2).Render(descText)

	// ── 3. Atalhos (sempre presentes, linha única) ─────────────────────────────
	shortcutHelp := fmt.Sprintf("%s %s   %s %s",
		theme.KeyShortcut.Render("[q] / [Ctrl+C]"),
		theme.KeyDesc.Render("Sair"),
		theme.KeyShortcut.Render("[r]"),
		theme.KeyDesc.Render("Tentar agora"),
	)

	// ── 4. Indicador de tentativas ────────────────────────────────────────────
	var attemptText string
	if attempts > 0 {
		attemptText = fmt.Sprintf("⏳ Aguardando... (tentativa %d)", attempts)
	} else {
		attemptText = "⏳ Aguardando o serviço ser iniciado..."
	}
	waitingMsg := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorAccent).Render(attemptText)

	// ── 5. Comandos sugeridos (omitidos se altura insuficiente) ───────────────
	var cmdBlock string
	if height >= 18 {
		cmdBoxWidth := boxWidth - 4
		if cmdBoxWidth < 20 {
			cmdBoxWidth = 20
		}

		cmdTitle := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorAccent).Render("Como iniciar o serviço:")

		lines := []string{cmdTitle, ""}

		if width >= 50 {
			lines = append(lines,
				lipgloss.NewStyle().Foreground(theme.ColorTextMuted).Render("• systemd (Ubuntu, Debian, Arch...):"),
				lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("  $ sudo systemctl start tailscaled"),
				"",
				lipgloss.NewStyle().Foreground(theme.ColorTextMuted).Render("• OpenRC (Alpine, Void Linux):"),
				lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("  $ sudo rc-service tailscaled start"),
				"",
				lipgloss.NewStyle().Foreground(theme.ColorTextMuted).Render("• Manual:"),
				lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("  $ sudo tailscaled"),
			)
		} else {
			// Tela estreita: comandos curtos
			lines = append(lines,
				lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo systemctl start tailscaled"),
				lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo rc-service tailscaled start"),
				lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo tailscaled"),
			)
		}

		cmdContent := lipgloss.JoinVertical(lipgloss.Left, lines...)
		cmdBlock = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.ColorBorder).
			Padding(0, 1).
			Width(cmdBoxWidth).
			Render(cmdContent)
	}

	// ── 6. Erro técnico truncado ──────────────────────────────────────────────
	var errSection string
	if lastError != "" && height >= 14 {
		cleanedErr := strings.ReplaceAll(lastError, "\n", " ")
		maxLen := boxWidth - 8
		if maxLen > 10 && len(cleanedErr) > maxLen {
			cleanedErr = cleanedErr[:maxLen] + "…"
		}
		errSection = lipgloss.NewStyle().Foreground(theme.ColorTextDim).
			Width(boxWidth - 2).Render("Log: " + cleanedErr)
	}

	// ── 7. Composição dinâmica dos elementos ──────────────────────────────────
	elements := []string{title, "", desc}
	if cmdBlock != "" {
		elements = append(elements, "", cmdBlock)
	}
	elements = append(elements, "", waitingMsg, "", shortcutHelp)
	if errSection != "" {
		elements = append(elements, "", errSection)
	}

	modalContent := lipgloss.JoinVertical(lipgloss.Left, elements...)

	// Borda externa: sem padding vertical extra em telas baixas
	paddingV := 1
	if height < 20 {
		paddingV = 0
	}

	mainBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.ColorWarning).
		Background(theme.ColorPanelBg).
		Padding(paddingV, 2).
		Width(boxWidth).
		MaxWidth(width).
		Render(modalContent)

	// Se a caixa renderizada for mais alta que a tela, usa Place sem centralizar
	// verticalmente para evitar corte — ancora no topo com 1 linha de margem.
	boxH := lipgloss.Height(mainBox)
	if boxH >= height {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Top, mainBox)
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, mainBox)
}
