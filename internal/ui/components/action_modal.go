package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"tuiscale/internal/tailscale"
	"tuiscale/internal/ui/theme"
)

// ExitNodeOptionItem representa uma opção de exit node no modal.
type ExitNodeOptionItem struct {
	Name     string
	HostName string
	IP       string
	Active   bool
}

// RenderExitNodeModal desenha o modal para seleção de Exit Node.
func RenderExitNodeModal(items []ExitNodeOptionItem, selectedIdx int, width, height int) string {
	modalWidth := 60
	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	title := theme.ModalTitle.Render("SELEÇÃO DE NÓ DE SAÍDA (EXIT NODE)")

	var rows []string
	for i, item := range items {
		isSelected := (i == selectedIdx)
		prefix := "  "
		if isSelected {
			prefix = "> "
		}

		activeMarker := ""
		if item.Active {
			activeMarker = " [EM USO]"
		}

		label := item.Name + activeMarker
		if isSelected {
			rows = append(rows, theme.TableRowSelected.Render(fmt.Sprintf("%s%-28s %s", prefix, label, item.IP)))
		} else {
			rows = append(rows, theme.TableRow.Render(fmt.Sprintf("%s%-28s %s", prefix, label, item.IP)))
		}
	}

	hint := lipgloss.NewStyle().Foreground(theme.ColorTextDim).Render("\n[↑/↓] Navegar   [Enter] Aplicar   [Esc] Cancelar")
	content := fmt.Sprintf("%s\n\n%s\n%s", title, strings.Join(rows, "\n"), hint)

	box := theme.ModalBox.Width(modalWidth).Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

// BuildExitNodeList monta a lista de exit nodes disponíveis a partir do status.
func BuildExitNodeList(status *tailscale.Status) []ExitNodeOptionItem {
	var items []ExitNodeOptionItem

	// Opção para desabilitar
	items = append(items, ExitNodeOptionItem{
		Name:     "Nenhum (Desativar Exit Node)",
		HostName: "",
		IP:       "",
		Active:   false,
	})

	if status == nil {
		return items
	}

	for _, p := range status.Peer {
		if p == nil {
			continue
		}
		if p.ExitNode || p.ExitNodeOption {
			ip := ""
			if len(p.TailscaleIPs) > 0 {
				ip = p.TailscaleIPs[0]
			}
			items = append(items, ExitNodeOptionItem{
				Name:     p.HostName,
				HostName: p.HostName,
				IP:       ip,
				Active:   p.ExitNode,
			})
		}
	}

	return items
}
