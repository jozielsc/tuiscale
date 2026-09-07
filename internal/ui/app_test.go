package ui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"tailscaletui/internal/tailscale"
)

func TestAppModelViews(t *testing.T) {
	client := tailscale.NewClient("")
	app := NewAppModel(client, 2*time.Second)

	// Simular evento de redimensionamento do terminal (120 colunas x 35 linhas)
	_, _ = app.Update(tea.WindowSizeMsg{Width: 120, Height: 35})

	// Simular status com peers
	status := &tailscale.Status{
		BackendState: "Running",
		TUN:          true,
		HaveNodeKey:  true,
		Self: &tailscale.PeerStatus{
			HostName:     "local-box",
			OS:           "linux",
			TailscaleIPs: []string{"100.113.167.75"},
			Active:       true,
		},
		CurrentTailnet: &tailscale.Tailnet{
			Name: "user@example.com",
		},
		Peer: map[string]*tailscale.PeerStatus{
			"p1": {
				HostName:     "peer-remote",
				OS:           "linux",
				TailscaleIPs: []string{"100.117.21.64"},
				Active:       true,
				CurAddr:      "179.155.132.31:51022",
				RxBytes:      1024000,
				TxBytes:      512000,
			},
		},
	}

	_, _ = app.Update(statusMsg{status: status, err: nil})

	// 1. Validar renderização da Aba 0 (Peers)
	view0 := app.View()
	if !strings.Contains(view0, "TailscaleTUI") {
		t.Errorf("View should contain app title 'TailscaleTUI'")
	}
	if !strings.Contains(view0, "CONECTADO") {
		t.Errorf("View should contain 'CONECTADO'")
	}
	if !strings.Contains(view0, "peer-remote") {
		t.Errorf("View should contain peer name 'peer-remote'")
	}

	// 2. Validar troca para Aba 1 (Netcheck)
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	view1 := app.View()
	if !strings.Contains(strings.ToLower(view1), "netcheck") {
		t.Errorf("Aba 1 should contain netcheck references, got: %s", view1)
	}

	// 3. Validar troca para Aba 2 (Métricas)
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	view2 := app.View()
	if !strings.Contains(view2, "MÉTRICAS") {
		t.Errorf("Aba 2 should contain metrics references")
	}

	// 4. Validar Modal de Ajuda
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	viewHelp := app.View()
	if !strings.Contains(viewHelp, "ATALHOS DO TECLADO") {
		t.Errorf("Help view should contain shortcuts guide")
	}
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyEsc})

	// 5. Validar Modal de Exit Node
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	viewExit := app.View()
	if !strings.Contains(viewExit, "EXIT NODE") {
		t.Errorf("Exit node view should contain EXIT NODE title")
	}
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyEsc})

	// 6. Validar Modal de Ping
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	viewPing := app.View()
	if !strings.Contains(viewPing, "PING") {
		t.Errorf("Ping view should contain PING title")
	}
}
