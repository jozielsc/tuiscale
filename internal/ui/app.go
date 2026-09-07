package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tailscaletui/internal/tailscale"
	"tailscaletui/internal/ui/clipboard"
	"tailscaletui/internal/ui/components"
	"tailscaletui/internal/ui/theme"
)

// Tipos de mensagens internas do Bubble Tea
type statusMsg struct {
	status *tailscale.Status
	err    error
}

type netcheckMsg struct {
	report *tailscale.NetcheckReport
	err    error
}

type pingResultMsg struct {
	lines []string
	err   error
}

type actionResultMsg struct {
	message string
	err     error
}

type clearToastMsg struct{}

// AppModel é o modelo principal da interface Bubble Tea.
type AppModel struct {
	client          *tailscale.Client
	refreshInterval time.Duration
	status          *tailscale.Status
	netcheckReport  *tailscale.NetcheckReport
	peerTable       *components.PeerTableModel

	width  int
	height int

	activeTab int // 0 = Peers, 1 = Rede/Netcheck, 2 = Métricas

	// Filtro de busca
	filterMode bool
	filterText string

	// Modais e estados interativos
	showHelp           bool
	showPeerDetail     bool
	showExitNodeModal  bool
	exitNodeItems      []components.ExitNodeOptionItem
	exitNodeCursor     int
	showPingModal      bool
	pingTargetHost     string
	pingTargetIP       string
	pingLines          []string
	isPingRunning      bool
	isNetcheckRunning  bool

	// Mensagens temporárias (toasts)
	toastMessage string
	toastIsError bool
}

// NewAppModel instancia o modelo do TailscaleTUI.
func NewAppModel(client *tailscale.Client, refresh time.Duration) *AppModel {
	if refresh <= 0 {
		refresh = 2 * time.Second
	}
	return &AppModel{
		client:          client,
		refreshInterval: refresh,
		peerTable:       components.NewPeerTableModel(),
		activeTab:       0,
	}
}

// Init inicializa os primeiros comandos assíncronos.
func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchStatusCmd(),
		m.tickCmd(),
	)
}

func (m *AppModel) tickCmd() tea.Cmd {
	return tea.Tick(m.refreshInterval, func(t time.Time) tea.Msg {
		return m.fetchStatusCmd()()
	})
}

func (m *AppModel) fetchStatusCmd() tea.Cmd {
	return func() tea.Msg {
		st, err := m.client.GetStatus(context.Background())
		return statusMsg{status: st, err: err}
	}
}

func (m *AppModel) runNetcheckCmd() tea.Cmd {
	m.isNetcheckRunning = true
	return func() tea.Msg {
		report, err := m.client.RunNetcheck(context.Background())
		return netcheckMsg{report: report, err: err}
	}
}

func (m *AppModel) pingPeerCmd(targetHost, targetIP string) tea.Cmd {
	m.isPingRunning = true
	m.pingLines = nil
	m.pingTargetHost = targetHost
	m.pingTargetIP = targetIP
	m.showPingModal = true
	return func() tea.Msg {
		lines, err := m.client.PingPeer(context.Background(), targetIP, 4)
		return pingResultMsg{lines: lines, err: err}
	}
}

func (m *AppModel) connectCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.client.Connect(context.Background())
		if err != nil {
			return actionResultMsg{message: "Falha ao conectar: " + err.Error(), err: err}
		}
		return actionResultMsg{message: "Tailscale conectado com sucesso!", err: nil}
	}
}

func (m *AppModel) disconnectCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.client.Disconnect(context.Background())
		if err != nil {
			return actionResultMsg{message: "Falha ao desconectar: " + err.Error(), err: err}
		}
		return actionResultMsg{message: "Tailscale desconectado!", err: nil}
	}
}

func (m *AppModel) setExitNodeCmd(node string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.SetExitNode(context.Background(), node)
		if err != nil {
			return actionResultMsg{message: "Falha ao alterar Exit Node: " + err.Error(), err: err}
		}
		msg := "Exit Node desativado"
		if node != "" {
			msg = "Exit Node configurado para: " + node
		}
		return actionResultMsg{message: msg, err: nil}
	}
}

func (m *AppModel) setToast(msg string, isError bool) tea.Cmd {
	m.toastMessage = msg
	m.toastIsError = isError
	return tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
		return clearToastMsg{}
	})
}

// Update processa eventos e mensagens.
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case statusMsg:
		if msg.err != nil {
			m.toastMessage = "Aviso: " + msg.err.Error()
			m.toastIsError = true
		} else {
			m.status = msg.status
			m.peerTable.UpdatePeers(m.status, m.filterText)
		}
		cmds = append(cmds, m.tickCmd())
		return m, tea.Batch(cmds...)

	case netcheckMsg:
		m.isNetcheckRunning = false
		if msg.err != nil {
			cmds = append(cmds, m.setToast("Erro no netcheck: "+msg.err.Error(), true))
		} else {
			m.netcheckReport = msg.report
			cmds = append(cmds, m.setToast("Netcheck concluído com sucesso!", false))
		}
		return m, tea.Batch(cmds...)

	case pingResultMsg:
		m.isPingRunning = false
		if msg.err != nil && len(msg.lines) == 0 {
			m.pingLines = []string{"Erro ao executar ping: " + msg.err.Error()}
		} else {
			m.pingLines = msg.lines
		}
		return m, nil

	case actionResultMsg:
		cmds = append(cmds, m.setToast(msg.message, msg.err != nil))
		cmds = append(cmds, m.fetchStatusCmd())
		return m, tea.Batch(cmds...)

	case clearToastMsg:
		m.toastMessage = ""
		return m, nil

	case tea.KeyMsg:
		// 1. Processamento quando o modo de busca está ativo
		if m.filterMode {
			switch msg.String() {
			case "enter", "esc":
				m.filterMode = false
			case "backspace":
				if len(m.filterText) > 0 {
					m.filterText = m.filterText[:len(m.filterText)-1]
					m.peerTable.UpdatePeers(m.status, m.filterText)
				}
			default:
				if len(msg.String()) == 1 {
					m.filterText += msg.String()
					m.peerTable.UpdatePeers(m.status, m.filterText)
				}
			}
			return m, nil
		}

		// 2. Modais ativos (fecham com Esc)
		if m.showHelp {
			if msg.String() == "esc" || msg.String() == "?" || msg.String() == "q" {
				m.showHelp = false
			}
			return m, nil
		}

		if m.showPeerDetail {
			if msg.String() == "esc" || msg.String() == "enter" || msg.String() == "q" {
				m.showPeerDetail = false
			}
			return m, nil
		}

		if m.showPingModal {
			if msg.String() == "esc" || msg.String() == "enter" || msg.String() == "q" {
				m.showPingModal = false
			}
			return m, nil
		}

		if m.showExitNodeModal {
			switch msg.String() {
			case "esc", "q":
				m.showExitNodeModal = false
			case "up", "k":
				if m.exitNodeCursor > 0 {
					m.exitNodeCursor--
				}
			case "down", "j":
				if m.exitNodeCursor < len(m.exitNodeItems)-1 {
					m.exitNodeCursor++
				}
			case "enter":
				selected := m.exitNodeItems[m.exitNodeCursor]
				m.showExitNodeModal = false
				return m, m.setExitNodeCmd(selected.HostName)
			}
			return m, nil
		}

		// 3. Atalhos globais e de navegação normal
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "?":
			m.showHelp = true
			return m, nil

		case "tab", "l", "right":
			m.activeTab = (m.activeTab + 1) % 3
			return m, nil

		case "shift+tab", "h", "left":
			m.activeTab = (m.activeTab + 2) % 3
			return m, nil

		case "1":
			m.activeTab = 0
			return m, nil

		case "2":
			m.activeTab = 1
			if m.netcheckReport == nil && !m.isNetcheckRunning {
				return m, m.runNetcheckCmd()
			}
			return m, nil

		case "3":
			m.activeTab = 2
			return m, nil

		case "/":
			m.filterMode = true
			return m, nil

		case "esc":
			if m.filterText != "" {
				m.filterText = ""
				m.peerTable.UpdatePeers(m.status, "")
			}
			return m, nil

		case "c":
			return m, m.connectCmd()

		case "d":
			return m, m.disconnectCmd()

		case "r":
			return m, m.fetchStatusCmd()

		case "n":
			return m, m.runNetcheckCmd()

		case "e":
			m.exitNodeItems = components.BuildExitNodeList(m.status)
			m.exitNodeCursor = 0
			m.showExitNodeModal = true
			return m, nil

		case "y":
			selected := m.peerTable.SelectedPeer()
			if selected != nil && len(selected.TailscaleIPs) > 0 {
				_ = clipboard.Copy(selected.TailscaleIPs[0])
				return m, m.setToast(fmt.Sprintf("IP %s copiado para a área de transferência!", selected.TailscaleIPs[0]), false)
			}

		case "Y":
			selected := m.peerTable.SelectedPeer()
			if selected != nil && selected.DNSName != "" {
				dns := strings.TrimSuffix(selected.DNSName, ".")
				_ = clipboard.Copy(dns)
				return m, m.setToast(fmt.Sprintf("DNS %s copiado para a área de transferência!", dns), false)
			}

		case "p":
			selected := m.peerTable.SelectedPeer()
			if selected != nil && len(selected.TailscaleIPs) > 0 {
				return m, m.pingPeerCmd(selected.HostName, selected.TailscaleIPs[0])
			}

		case "enter":
			if m.activeTab == 0 && m.peerTable.SelectedPeer() != nil {
				m.showPeerDetail = true
				return m, nil
			}

		case "up", "k":
			if m.activeTab == 0 {
				m.peerTable.MoveUp()
			}

		case "down", "j":
			if m.activeTab == 0 {
				m.peerTable.MoveDown()
			}

		case "g", "home":
			if m.activeTab == 0 {
				m.peerTable.MoveTop()
			}

		case "G", "end":
			if m.activeTab == 0 {
				m.peerTable.MoveBottom()
			}
		}
	}

	return m, nil
}

// View renderiza a interface no terminal.
func (m *AppModel) View() string {
	if m.width < 50 || m.height < 15 {
		return "Janela do terminal muito pequena para renderizar o TailscaleTUI. Redimensione a janela."
	}

	// 1. Header
	headerView := components.RenderHeader(m.status, m.client.GetSpeedTracker(), m.width)

	// 2. Tab Bar
	tabPeers := theme.TabInactive.Render(" [1] Dispositivos / Peers ")
	tabNet := theme.TabInactive.Render(" [2] Rede & Netcheck ")
	tabMetrics := theme.TabInactive.Render(" [3] Tráfego & Métricas ")

	switch m.activeTab {
	case 0:
		tabPeers = theme.TabActive.Render(" [1] Dispositivos / Peers ")
	case 1:
		tabNet = theme.TabActive.Render(" [2] Rede & Netcheck ")
	case 2:
		tabMetrics = theme.TabActive.Render(" [3] Tráfego & Métricas ")
	}

	tabBar := theme.TabBar.Render(lipgloss.JoinHorizontal(lipgloss.Top, tabPeers, tabNet, tabMetrics))

	// 3. Cálculo do espaço para o corpo principal
	headerH := lipgloss.Height(headerView)
	tabH := lipgloss.Height(tabBar)
	footerH := 1
	bodyH := m.height - headerH - tabH - footerH - 2
	if bodyH < 5 {
		bodyH = 5
	}

	// 4. Renderização da Aba Ativa
	var bodyView string
	switch m.activeTab {
	case 0:
		// Em telas largas (> 110 colunas), layout lado a lado: Lista à esquerda e Detalhes à direita!
		if m.width > 115 {
			leftW := (m.width * 6) / 10
			rightW := m.width - leftW - 3
			leftContent := m.peerTable.Render(leftW, bodyH)
			rightContent := components.RenderPeerDetail(m.peerTable.SelectedPeer(), rightW, bodyH)
			bodyView = lipgloss.JoinHorizontal(lipgloss.Top, leftContent, " ", rightContent)
		} else {
			bodyView = m.peerTable.Render(m.width-2, bodyH)
		}

	case 1:
		bodyView = components.RenderNetcheckView(m.netcheckReport, m.isNetcheckRunning, m.width-2, bodyH)

	case 2:
		bodyView = components.RenderMetricsView(m.status, m.client.GetSpeedTracker(), m.width-2, bodyH)
	}

	// 5. Footer
	footerView := components.RenderFooter(m.toastMessage, m.toastIsError, m.filterMode, m.filterText, m.width)

	// Composição da tela base
	mainView := lipgloss.JoinVertical(lipgloss.Left,
		headerView,
		tabBar,
		bodyView,
		footerView,
	)

	// 6. Sobreposição de Modais se houver algum ativo
	if m.showHelp {
		return components.RenderHelpModal(m.width, m.height)
	}
	if m.showPeerDetail && m.width <= 115 {
		// Em telas compactas onde o detalhe não está lado a lado
		box := components.RenderPeerDetail(m.peerTable.SelectedPeer(), m.width-6, m.height-6)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
	if m.showPingModal {
		return components.RenderPingModal(m.pingTargetHost, m.pingTargetIP, m.pingLines, m.isPingRunning, m.width, m.height)
	}
	if m.showExitNodeModal {
		return components.RenderExitNodeModal(m.exitNodeItems, m.exitNodeCursor, m.width, m.height)
	}

	return mainView
}
