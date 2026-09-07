package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"tuiscale/internal/ui/theme"
)

// RenderDaemonWaitView desenha a tela amigável informando que o daemon tailscaled está inativo e aguardando inicialização.
func RenderDaemonWaitView(width, height int, lastError string, attempts int) string {
	boxWidth := 74
	if width-4 < boxWidth {
		boxWidth = width - 4
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	// 1. Título e Badge
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
		descText = "O TUIScale não conseguiu se comunicar com o daemon local ('tailscaled').\nO serviço parece estar desligado ou ainda não foi inicializado neste sistema."
	}

	title := lipgloss.JoinHorizontal(lipgloss.Center,
		theme.TitleStyle.Render("TUIScale"),
		"  ",
		badge,
	)

	// 2. Descrição Amigável
	desc := lipgloss.NewStyle().
		Foreground(theme.ColorText).
		Render(descText)

	// 3. Caixa de Comandos Sugeridos
	cmdBoxWidth := boxWidth - 6
	if cmdBoxWidth < 34 {
		cmdBoxWidth = 34
	}

	cmdTitle := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorAccent).Render("Comandos sugeridos para iniciar o serviço:")

	systemdLabel := lipgloss.NewStyle().Foreground(theme.ColorTextMuted).Render("• Distribuições com systemd (Ubuntu, Debian, Fedora, Arch, etc.):")
	systemdCmd := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("  $ sudo systemctl start tailscaled")

	openrcLabel := lipgloss.NewStyle().Foreground(theme.ColorTextMuted).Render("• Distribuições com OpenRC (Alpine, Void Linux):")
	openrcCmd := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("  $ sudo rc-service tailscaled start")

	manualLabel := lipgloss.NewStyle().Foreground(theme.ColorTextMuted).Render("• Execução manual direta:")
	manualCmd := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("  $ sudo tailscaled")

	cmdContent := lipgloss.JoinVertical(lipgloss.Left,
		cmdTitle,
		"",
		systemdLabel,
		systemdCmd,
		"",
		openrcLabel,
		openrcCmd,
		"",
		manualLabel,
		manualCmd,
	)

	cmdBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.ColorBorder).
		Padding(0, 1).
		Width(cmdBoxWidth).
		Render(cmdContent)

	// 4. Indicador de Espera / Polling
	var attemptText string
	if attempts > 0 {
		attemptText = fmt.Sprintf(" (tentativa %d - checando a cada 2s)", attempts)
	} else {
		attemptText = " (checando a cada 2s)"
	}

	waitingMsg := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.ColorAccent).
		Render("⏳ Aguardando o serviço ser iniciado..." + attemptText)

	// 5. Atalhos Rápidos
	shortcutHelp := fmt.Sprintf("%s %s      %s %s",
		theme.KeyShortcut.Render("[q] / [Ctrl+C]"),
		theme.KeyDesc.Render("Sair"),
		theme.KeyShortcut.Render("[r]"),
		theme.KeyDesc.Render("Tentar conectar agora"),
	)

	// 6. Detalhes do erro técnico (truncado se necessário)
	var errSection string
	if lastError != "" {
		cleanedErr := strings.ReplaceAll(lastError, "\n", " ")
		maxLen := boxWidth - 14
		if maxLen > 10 && len(cleanedErr) > maxLen {
			cleanedErr = cleanedErr[:maxLen] + "..."
		}
		errSection = lipgloss.NewStyle().
			Foreground(theme.ColorTextDim).
			Render("Log: " + cleanedErr)
	}

	// Compor elementos verticais da caixa
	elements := []string{
		title,
		"",
		desc,
		"",
		cmdBox,
		"",
		waitingMsg,
		"",
		shortcutHelp,
	}

	if errSection != "" {
		elements = append(elements, "", errSection)
	}

	modalContent := lipgloss.JoinVertical(lipgloss.Left, elements...)

	mainBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.ColorWarning).
		Background(theme.ColorPanelBg).
		Padding(1, 2).
		Width(boxWidth).
		Render(modalContent)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, mainBox)
}
