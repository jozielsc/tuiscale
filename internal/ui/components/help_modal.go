package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"tailscaletui/internal/ui/theme"
)

// RenderHelpModal desenha a janela modal de ajuda com todos os atalhos disponíveis.
func RenderHelpModal(width, height int) string {
	modalWidth := 65
	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	title := theme.ModalTitle.Render("ATALHOS DO TECLADO - TAILSCALETUI")

	type shortcut struct {
		key  string
		desc string
	}

	renderSection := func(secTitle string, items []shortcut) string {
		out := []string{lipgloss.NewStyle().Bold(true).Foreground(theme.ColorAccent).Render(secTitle)}
		for _, it := range items {
			line := fmt.Sprintf("  %-14s %s",
				theme.KeyShortcut.Render(it.key),
				theme.KeyDesc.Render(it.desc),
			)
			out = append(out, line)
		}
		return strings.Join(out, "\n")
	}

	secNav := renderSection("Navegação:", []shortcut{
		{"↑ / k, ↓ / j", "Navegar na lista de dispositivos"},
		{"PgUp / PgDn", "Subir ou descer uma página"},
		{"Tab / 1, 2, 3", "Alternar entre abas (Peers / Rede / Métricas)"},
		{"Enter", "Visualizar detalhes aprofundados do dispositivo"},
	})

	secActions := renderSection("Ações do Tailscale:", []shortcut{
		{"c", "Conectar Tailscale (tailscale up)"},
		{"d", "Desconectar Tailscale (tailscale down)"},
		{"p", "Disparar teste de ping no peer selecionado"},
		{"e", "Selecionar ou desativar Nó de Saída (Exit Node)"},
		{"y", "Copiar IP Tailscale do peer selecionado"},
		{"Y", "Copiar nome MagicDNS do peer selecionado"},
	})

	secTools := renderSection("Busca e Diagnósticos:", []shortcut{
		{"/", "Filtrar dispositivos por nome, IP, SO ou usuário"},
		{"n", "Executar diagnóstico completo de rede (netcheck)"},
		{"r", "Forçar atualização de status"},
		{"Esc", "Limpar busca ou fechar qualquer janela modal"},
		{"?", "Abrir ou fechar esta janela de ajuda"},
		{"q / Ctrl+C", "Sair do aplicativo"},
	})

	hint := lipgloss.NewStyle().Foreground(theme.ColorTextDim).Render("\nPressione [Esc] ou [?] para fechar a ajuda")
	content := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n%s", title, secNav, secActions, secTools, hint)

	box := theme.ModalBox.Width(modalWidth).Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}
