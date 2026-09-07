package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Cores do Dark Mode
var (
	ColorBgDark    = lipgloss.Color("#12141A")
	ColorPanelBg   = lipgloss.Color("#181B20")
	ColorBorder    = lipgloss.Color("#2D333F")
	ColorBorderDim = lipgloss.Color("#222731")
	ColorHighlight = lipgloss.Color("#3B82F6") // Tailscale Blue
	ColorAccent    = lipgloss.Color("#38BDF8") // Sky Blue / Cyan
	ColorText      = lipgloss.Color("#F1F5F9") // Light gray/white
	ColorTextDim   = lipgloss.Color("#64748B") // Slate gray
	ColorTextMuted = lipgloss.Color("#94A3B8") // Medium gray
	ColorSuccess   = lipgloss.Color("#22C55E") // Green
	ColorWarning   = lipgloss.Color("#EAB308") // Yellow/Amber
	ColorDanger    = lipgloss.Color("#EF4444") // Red
	ColorPurple    = lipgloss.Color("#A855F7") // Purple
	ColorRowSelect = lipgloss.Color("#263043") // Row highlight background
)

// Estilos de UI
var (
	// Base
	BaseStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	// Header
	HeaderBox = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorHighlight).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHighlight)

	SubTitleStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted)

	// Status Badges
	BadgeRunning = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorSuccess).
			Padding(0, 1)

	BadgeStopped = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorDanger).
			Padding(0, 1)

	BadgeStarting = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorWarning).
			Padding(0, 1)

	// Connection Mode Badges
	BadgeDirect = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	BadgeRelay = lipgloss.NewStyle().
			Foreground(ColorWarning)

	BadgeOffline = lipgloss.NewStyle().
			Foreground(ColorTextDim)

	BadgeExitNode = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPurple).
			Padding(0, 1)

	// Tabs
	TabBar = lipgloss.NewStyle().
		MarginBottom(0)

	TabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHighlight).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorHighlight).
			Padding(0, 2)

	TabInactive = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Padding(0, 2)

	// Tables and Lists
	TableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	TableRow = lipgloss.NewStyle().
			Padding(0, 1)

	TableRowSelected = lipgloss.NewStyle().
				Background(ColorRowSelect).
				Bold(true).
				Foreground(ColorText).
				Padding(0, 1)

	// Panels
	Panel = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	PanelActive = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorHighlight).
			Padding(0, 1)

	// Speed meters
	SpeedRxStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	SpeedTxStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	// Modais
	ModalBox = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(ColorHighlight).
			Background(ColorPanelBg).
			Padding(1, 2)

	ModalTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			MarginBottom(1)

	// Footer / Statusbar
	FooterBar = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorBorderDim).
			Padding(0, 1)

	KeyShortcut = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHighlight)

	KeyDesc = lipgloss.NewStyle().
		Foreground(ColorTextMuted)

	ToastSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	ToastError = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)
)

// FormatOSIcon retorna uma tag ou ícone formatado para o sistema operacional.
func FormatOSIcon(osName string) string {
	switch osName {
	case "linux":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#E2B340")).Render("🐧 Linux")
	case "windows":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00A4EF")).Render("🪟 Windows")
	case "darwin", "macos", "mac":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F5F5F7")).Render("🍎 macOS")
	case "android":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#3DDC84")).Render("🤖 Android")
	case "ios":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#A2AAAD")).Render("📱 iOS")
	case "freebsd", "openbsd":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#AB2B28")).Render("😈 BSD")
	default:
		return lipgloss.NewStyle().Foreground(ColorTextDim).Render("💻 " + osName)
	}
}
