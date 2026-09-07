package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap define as teclas de atalho da aplicação.
type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	NextTab    key.Binding
	PrevTab    key.Binding
	Tab1       key.Binding
	Tab2       key.Binding
	Tab3       key.Binding
	Connect    key.Binding
	Disconnect key.Binding
	Ping       key.Binding
	ExitNode   key.Binding
	CopyIP     key.Binding
	CopyDNS    key.Binding
	Filter     key.Binding
	Refresh    key.Binding
	Enter      key.Binding
	Escape     key.Binding
	Help       key.Binding
	Quit       key.Binding
}

// DefaultKeyMap retorna o mapeamento padrão de teclas de atalho.
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "subir"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "descer"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "página anterior"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "próxima página"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab", "l", "right"),
		key.WithHelp("tab/→", "próxima aba"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab", "h", "left"),
		key.WithHelp("shift+tab/←", "aba anterior"),
	),
	Tab1: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "aba peers"),
	),
	Tab2: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "aba rede/netcheck"),
	),
	Tab3: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "aba métricas"),
	),
	Connect: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "conectar"),
	),
	Disconnect: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "desconectar"),
	),
	Ping: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "ping peer"),
	),
	ExitNode: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "exit node"),
	),
	CopyIP: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "copiar IP"),
	),
	CopyDNS: key.NewBinding(
		key.WithKeys("Y"),
		key.WithHelp("Y", "copiar DNS"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filtrar peers"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "atualizar"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "detalhes"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "fechar/cancelar"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "ajuda"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "sair"),
	),
}
