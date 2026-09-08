package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/tailscale"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// PeerTableModel gerencia o estado da tabela de peers.
type PeerTableModel struct {
	cursor       int
	scrollOffset int
	filter       string
	peers        []*tailscale.PeerStatus
}

// NewPeerTableModel cria uma nova instância de PeerTableModel.
func NewPeerTableModel() *PeerTableModel {
	return &PeerTableModel{
		cursor: 0,
		peers:  make([]*tailscale.PeerStatus, 0),
	}
}

// UpdatePeers atualiza a lista de peers aplicando ordenação e filtro.
func (m *PeerTableModel) UpdatePeers(status *tailscale.Status, filter string) {
	m.filter = filter
	if status == nil {
		m.peers = nil
		return
	}

	var list []*tailscale.PeerStatus

	// Adicionar nó local se existir
	if status.Self != nil {
		selfCopy := *status.Self
		selfCopy.HostName = selfCopy.HostName + " (esta máquina)"
		list = append(list, &selfCopy)
	}

	// Adicionar peers remotos
	for _, p := range status.Peer {
		if p != nil {
			list = append(list, p)
		}
	}

	// Aplicar filtro de busca se houver
	var filtered []*tailscale.PeerStatus
	filterLower := strings.ToLower(filter)
	for _, p := range list {
		if filter == "" {
			filtered = append(filtered, p)
			continue
		}
		match := strings.Contains(strings.ToLower(p.HostName), filterLower) ||
			strings.Contains(strings.ToLower(p.OS), filterLower) ||
			strings.Contains(strings.ToLower(p.UserLogin), filterLower)
		if !match {
			for _, ip := range p.TailscaleIPs {
				if strings.Contains(ip, filterLower) {
					match = true
					break
				}
			}
		}
		if match {
			filtered = append(filtered, p)
		}
	}

	// Ordenação: 1º Active/Online, 2º Nome do host
	sort.SliceStable(filtered, func(i, j int) bool {
		pi, pj := filtered[i], filtered[j]
		if pi.Active != pj.Active {
			return pi.Active
		}
		if pi.Online != pj.Online {
			return pi.Online
		}
		return strings.ToLower(pi.HostName) < strings.ToLower(pj.HostName)
	})

	m.peers = filtered

	// Ajustar cursor para não ultrapassar a lista
	if m.cursor >= len(m.peers) {
		m.cursor = len(m.peers) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// MoveUp move o cursor uma linha para cima.
func (m *PeerTableModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown move o cursor uma linha para baixo.
func (m *PeerTableModel) MoveDown() {
	if m.cursor < len(m.peers)-1 {
		m.cursor++
	}
}

// MoveTop move o cursor para o topo.
func (m *PeerTableModel) MoveTop() {
	m.cursor = 0
}

// MoveBottom move o cursor para o final.
func (m *PeerTableModel) MoveBottom() {
	if len(m.peers) > 0 {
		m.cursor = len(m.peers) - 1
	}
}

// RecalculateScroll ajusta o scroll offset após mudança de tamanho de tela.
func (m *PeerTableModel) RecalculateScroll(visibleRows int) {
	if len(m.peers) == 0 {
		return
	}

	// Ajustar scroll offset para manter o cursor visível
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	} else if m.cursor >= m.scrollOffset+visibleRows {
		m.scrollOffset = m.cursor - visibleRows + 1
	}

	// Garantir que scrollOffset não seja negativo
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

// SelectedPeer retorna o peer atualmente selecionado.
func (m *PeerTableModel) SelectedPeer() *tailscale.PeerStatus {
	if len(m.peers) == 0 || m.cursor < 0 || m.cursor >= len(m.peers) {
		return nil
	}
	return m.peers[m.cursor]
}

// Render desenha a tabela de peers com suporte a scroll e dimensionamento.
func (m *PeerTableModel) Render(width, height int) string {
	if len(m.peers) == 0 {
		emptyMsg := "Nenhum dispositivo encontrado."
		if m.filter != "" {
			emptyMsg = fmt.Sprintf("Nenhum dispositivo encontrado para o filtro: '%s'", m.filter)
		}
		return lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(theme.ColorTextDim).
			Render(emptyMsg)
	}

	// Definição de larguras de coluna proporcionais
	colStatusW := 4
	colHostW := 22
	colOSW := 12
	colIPW := 16
	colConnW := 22
	colTrafficW := 18

	// Adaptar para larguras menores se necessário
	if width < 100 {
		colHostW = 16
		colConnW = 16
		colTrafficW = 14
	}

	header := fmt.Sprintf("  %-*s %-*s %-*s %-*s %-*s %-*s",
		colStatusW, "ST",
		colHostW, "DISPOSITIVO",
		colOSW, "SO",
		colIPW, "IP TAILSCALE",
		colConnW, "ROTA / CONEXÃO",
		colTrafficW, "TRÁFEGO (RX/TX)",
	)
	headerRendered := theme.TableHeader.Width(width).Render(header)

	// Quantidade de linhas visíveis
	visibleRows := height - 3 // Descontar cabeçalho da tabela e bordas
	if visibleRows < 1 {
		visibleRows = 1
	}

	// Ajuste do scroll offset para manter o cursor visível
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	} else if m.cursor >= m.scrollOffset+visibleRows {
		m.scrollOffset = m.cursor - visibleRows + 1
	}

	endIdx := m.scrollOffset + visibleRows
	if endIdx > len(m.peers) {
		endIdx = len(m.peers)
	}

	var rows []string
	for i := m.scrollOffset; i < endIdx; i++ {
		p := m.peers[i]
		isSelected := (i == m.cursor)

		// 1. Status icon
		var statusIcon string
		if p.Active {
			statusIcon = theme.BadgeDirect.Render("🟢")
		} else if p.Online {
			statusIcon = theme.BadgeRelay.Render("🔵")
		} else {
			statusIcon = theme.BadgeOffline.Render("⚪")
		}

		// 2. Hostname e ExitNode tag
		hostName := p.HostName
		if p.ExitNode {
			hostName += " [EXIT]"
		}
		if len(hostName) > colHostW-1 {
			hostName = hostName[:colHostW-2] + "…"
		}

		// 3. Sistema Operacional
		osName := theme.FormatOSIcon(strings.ToLower(p.OS))

		// 4. IP
		ip := "-"
		if len(p.TailscaleIPs) > 0 {
			ip = p.TailscaleIPs[0]
		}

		// 5. Conexão / Rota
		var connStr string
		if p.CurAddr != "" {
			connStr = "Direta: " + p.CurAddr
		} else if p.Relay != "" {
			connStr = "Relay DERP: " + p.Relay
		} else if p.Online {
			connStr = "Online (Idle)"
		} else {
			connStr = "Offline"
		}
		if len(connStr) > colConnW-1 {
			connStr = connStr[:colConnW-2] + "…"
		}

		// 6. Tráfego
		trafficStr := fmt.Sprintf("↓%s ↑%s",
			tailscale.FormatBytes(p.RxBytes),
			tailscale.FormatBytes(p.TxBytes),
		)

		line := fmt.Sprintf(" %s %-*s %-*s %-*s %-*s %-*s",
			statusIcon,
			colHostW, hostName,
			colOSW, osName,
			colIPW, ip,
			colConnW, connStr,
			colTrafficW, trafficStr,
		)

		if isSelected {
			rows = append(rows, theme.TableRowSelected.Width(width).Render(">"+line))
		} else {
			rows = append(rows, theme.TableRow.Width(width).Render(" "+line))
		}
	}

	content := headerRendered + "\n" + strings.Join(rows, "\n")
	return content
}
