package ui

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jozielsc/tuiscale/internal/tailscale"
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
	if !strings.Contains(view0, "TUIScale") {
		t.Errorf("View should contain app title 'TUIScale'")
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

func TestDaemonWaitViewAndTransitions(t *testing.T) {
	client := tailscale.NewClient("")
	app := NewAppModel(client, 2*time.Second)

	_, _ = app.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// 1. Verificar estado inicial de espera
	if !app.daemonWaiting {
		t.Errorf("app should initially be in daemonWaiting state")
	}

	viewInitial := app.View()
	if !strings.Contains(viewInitial, "TUIScale") {
		t.Errorf("initial view should contain 'TUIScale'")
	}

	// 2. Simular falha de conexão (daemon inativo)
	daemonErr := errors.New("failed to connect to local tailscaled: dial unix /var/run/tailscale/tailscaled.sock: connect: no such file or directory")
	_, _ = app.Update(statusMsg{status: nil, err: daemonErr})

	if !app.daemonWaiting {
		t.Errorf("app should remain in daemonWaiting state after connection error")
	}
	if app.daemonAttempts != 1 {
		t.Errorf("expected daemonAttempts to be 1, got %d", app.daemonAttempts)
	}

	viewWaiting := app.View()
	if !strings.Contains(viewWaiting, "SERVIÇO TAILSCALE INATIVO") {
		t.Errorf("waiting view should contain 'SERVIÇO TAILSCALE INATIVO'")
	}
	if !strings.Contains(viewWaiting, "sudo systemctl start tailscaled") {
		t.Errorf("waiting view should suggest systemctl start tailscaled")
	}
	if !strings.Contains(viewWaiting, "sudo rc-service tailscaled start") {
		t.Errorf("waiting view should suggest rc-service start")
	}
	if !strings.Contains(viewWaiting, "sudo tailscaled") {
		t.Errorf("waiting view should suggest manual execution sudo tailscaled")
	}
	if !strings.Contains(viewWaiting, "[q] / [Ctrl+C]") {
		t.Errorf("waiting view should show exit shortcuts")
	}

	// 3. Testar atalhos no modo de espera
	// Pressionar 'q' deve retornar tea.Quit
	_, cmdQ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmdQ == nil || cmdQ() != tea.Quit() {
		t.Errorf("expected 'q' to return tea.Quit while in daemonWaiting state")
	}

	// Pressionar 'ctrl+c' deve retornar tea.Quit
	_, cmdCtrlC := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmdCtrlC == nil || cmdCtrlC() != tea.Quit() {
		t.Errorf("expected ctrl+c to return tea.Quit while in daemonWaiting state")
	}

	// Pressionar 'r' deve disparar verificação imediata
	_, cmdR := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if cmdR == nil {
		t.Errorf("expected 'r' to return a command while in daemonWaiting state")
	}

	// Pressionar outras teclas como '2' não deve trocar de aba
	_, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	if app.activeTab != 0 {
		t.Errorf("expected activeTab to remain 0 while waiting for daemon, got %d", app.activeTab)
	}

	// 4. Simular recuperação do serviço Tailscale
	runningStatus := &tailscale.Status{
		BackendState: "Running",
		Self: &tailscale.PeerStatus{
			HostName:     "local-node",
			OS:           "linux",
			TailscaleIPs: []string{"100.64.0.1"},
			Active:       true,
		},
	}
	_, _ = app.Update(statusMsg{status: runningStatus, err: nil})

	if app.daemonWaiting {
		t.Errorf("daemonWaiting should be false after receiving successful statusMsg")
	}

	viewActive := app.View()
	if strings.Contains(viewActive, "SERVIÇO TAILSCALE INATIVO") {
		t.Errorf("active view should not show daemon waiting screen")
	}
	if !strings.Contains(viewActive, "Dispositivos / Peers") {
		t.Errorf("active view should show tabs")
	}
}


func TestSmallWindowGuard(t *testing.T) {
	client := tailscale.NewClient("")

	// 1. Janela pequena com daemon conectado → mensagem base
	app := NewAppModel(client, 2*time.Second)
	_, _ = app.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	app.daemonWaiting = false
	view := app.View()
	if !strings.Contains(view, "muito pequeno") {
		t.Errorf("small window with daemon up should show size warning, got: %q", view)
	}
	if strings.Contains(view, "tailscaled inativo") {
		t.Errorf("small window with daemon up should NOT mention daemon inactive")
	}

	// 2. Janela pequena com daemon INATIVO → inclui dica do daemon
	app2 := NewAppModel(client, 2*time.Second)
	_, _ = app2.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	// daemonWaiting já é true por padrão em NewAppModel
	view2 := app2.View()
	if !strings.Contains(view2, "tailscaled inativo") {
		t.Errorf("small window with daemon down should mention daemon inactive, got: %q", view2)
	}
	if !strings.Contains(view2, "[q]") {
		t.Errorf("small window with daemon down should show quit shortcut")
	}
	if !strings.Contains(view2, "[r]") {
		t.Errorf("small window with daemon down should show retry shortcut")
	}
}
