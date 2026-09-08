package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// boxFits retorna true se o conteúdo, ao ser envolvido na caixa com a largura
// e padding dados, cabe dentro de availHeight linhas de terminal.
func boxFits(content string, boxWidth, paddingV, availHeight int) bool {
	// bordas superior + inferior = 2, padding vertical = paddingV*2
	overhead := 2 + paddingV*2
	return lipgloss.Height(content)+overhead <= availHeight
}

// RenderDaemonWaitView desenha a tela informando que o daemon tailscaled está inativo.
// Usa construção em camadas progressivas: parte do núcleo obrigatório (título + shortcuts
// + indicador) e adiciona elementos opcionais apenas enquanto couberem na altura disponível.
// Isso garante que a tela nunca transborde independentemente do tamanho da janela.
func RenderDaemonWaitView(width, height int, lastError string, attempts int) string {
	boxWidth := width - 4
	if boxWidth > 74 {
		boxWidth = 74
	}
	if boxWidth < 20 {
		boxWidth = 20
	}

	// Ultra-compacto: abaixo de 30 colunas ou 6 linhas exibe texto simples
	if width < 30 || height < 6 {
		msg := "tailscaled inativo  [q] Sair  [r] Retry"
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(theme.ColorWarning).Render(msg))
	}

	// ── Peças do conteúdo ────────────────────────────────────────────────────

	var badge string
	if attempts == 0 && lastError == "" {
		badge = theme.BadgeStarting.Render("⏳ CONECTANDO AO DAEMON")
	} else {
		badge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(theme.ColorWarning).
			Padding(0, 1).
			Render("⚠️  SERVIÇO TAILSCALE INATIVO")
	}

	title := lipgloss.JoinHorizontal(lipgloss.Center,
		theme.TitleStyle.Render("TUIScale"), "  ", badge,
	)

	// Atalhos — SEMPRE incluídos (núcleo obrigatório junto ao título)
	shortcutHelp := fmt.Sprintf("%s %s   %s %s",
		theme.KeyShortcut.Render("[q] / [Ctrl+C]"),
		theme.KeyDesc.Render("Sair"),
		theme.KeyShortcut.Render("[r]"),
		theme.KeyDesc.Render("Tentar agora"),
	)

	// Indicador de tentativas — SEMPRE incluído
	var attemptText string
	if attempts > 0 {
		attemptText = fmt.Sprintf("⏳ Aguardando... (tentativa %d)", attempts)
	} else {
		attemptText = "⏳ Aguardando o serviço ser iniciado..."
	}
	waitingMsg := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorAccent).Render(attemptText)

	// Largura útil interna do box: boxWidth menos bordas (2) e padding horizontal (2*2=4).
	// Elementos internos devem respeitar esse limite para não sofrerem re-wrap.
	innerContentWidth := boxWidth - 6
	if innerContentWidth < 14 {
		innerContentWidth = 14
	}

	// Descrição — opcional
	var descText string
	if attempts == 0 && lastError == "" {
		descText = "Iniciando comunicação com o daemon local do Tailscale..."
	} else {
		descText = "Não foi possível se comunicar com o daemon local ('tailscaled').\nO serviço parece estar desligado ou ainda não foi inicializado."
	}
	desc := lipgloss.NewStyle().Foreground(theme.ColorText).Width(innerContentWidth).Render(descText)

	// Bloco de comandos — opcional
	cmdBoxWidth := innerContentWidth
	if cmdBoxWidth < 20 {
		cmdBoxWidth = 20
	}
	var cmdLines []string
	cmdLines = append(cmdLines,
		lipgloss.NewStyle().Bold(true).Foreground(theme.ColorAccent).Render("Como iniciar o serviço:"),
		"",
	)
	if width >= 50 {
		cmdLines = append(cmdLines,
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
		cmdLines = append(cmdLines,
			lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo systemctl start tailscaled"),
			lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo rc-service tailscaled start"),
			lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo tailscaled"),
		)
	}
	cmdBlock := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.ColorBorder).
		Padding(0, 1).
		Width(cmdBoxWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, cmdLines...))

	// Log de erro — opcional
	var errSection string
	if lastError != "" {
		cleanedErr := strings.ReplaceAll(lastError, "\n", " ")
		maxLen := innerContentWidth - 6
		if maxLen > 10 && len(cleanedErr) > maxLen {
			cleanedErr = cleanedErr[:maxLen] + "…"
		}
		errSection = lipgloss.NewStyle().Foreground(theme.ColorTextDim).
			Width(innerContentWidth).Render("Log: " + cleanedErr)
	}

	// ── Construção em camadas progressivas ────────────────────────────────────
	// Começa com o núcleo obrigatório e adiciona camadas opcionais apenas
	// enquanto o resultado couber na altura disponível da janela.

	// Padding vertical da borda: 0 em telas baixas para economizar linhas
	paddingV := 1
	if height < 20 {
		paddingV = 0
	}

	// Núcleo: título + indicador + shortcuts (sempre presente)
	core := lipgloss.JoinVertical(lipgloss.Left,
		title, "", waitingMsg, "", shortcutHelp,
	)

	// Camada 1: + descrição
	withDesc := lipgloss.JoinVertical(lipgloss.Left,
		title, "", desc, "", waitingMsg, "", shortcutHelp,
	)

	// Camada 2: + bloco de comandos
	withCmds := lipgloss.JoinVertical(lipgloss.Left,
		title, "", desc, "", cmdBlock, "", waitingMsg, "", shortcutHelp,
	)

	// Camada 2b: versão compacta do cmdBlock (só comandos, sem labels descritivas)
	cmdCompact := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.ColorBorder).
		Padding(0, 1).
		Width(cmdBoxWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo systemctl start tailscaled"),
			lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo rc-service tailscaled start"),
			lipgloss.NewStyle().Bold(true).Foreground(theme.ColorHighlight).Render("sudo tailscaled"),
		))
	withCmdsCompact := lipgloss.JoinVertical(lipgloss.Left,
		title, "", desc, "", cmdCompact, "", waitingMsg, "", shortcutHelp,
	)

	// Camada 3: + log de erro
	withAll := withCmds
	if errSection != "" {
		withAll = lipgloss.JoinVertical(lipgloss.Left,
			title, "", desc, "", cmdBlock, "", waitingMsg, "", shortcutHelp, "", errSection,
		)
	}

	// Escolhe a versão mais completa que caiba
	var finalContent string
	switch {
	case errSection != "" && boxFits(withAll, boxWidth, paddingV, height):
		finalContent = withAll
	case boxFits(withCmds, boxWidth, paddingV, height):
		finalContent = withCmds
	case boxFits(withCmdsCompact, boxWidth, paddingV, height):
		finalContent = withCmdsCompact
	case boxFits(withDesc, boxWidth, paddingV, height):
		finalContent = withDesc
	default:
		finalContent = core
	}

	mainBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.ColorWarning).
		Background(theme.ColorPanelBg).
		Padding(paddingV, 2).
		Width(boxWidth).
		Render(finalContent)

	// Fallback de segurança: se o box renderizado ainda exceder height,
	// renderiza apenas o núcleo (sem desc, sem cmds) que é garantidamente pequeno.
	if lipgloss.Height(mainBox) > height {
		mainBox = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.ColorWarning).
			Background(theme.ColorPanelBg).
			Padding(0, 2).
			Width(boxWidth).
			Render(core)
	}

	placed := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, mainBox)
	// Place retorna o conteúdo sem truncar quando ele é maior que height.
	// Nesse caso, retornamos direto o mainBox (que já cabe por construção).
	if lipgloss.Height(placed) > height {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Top, mainBox)
	}
	return placed
}
