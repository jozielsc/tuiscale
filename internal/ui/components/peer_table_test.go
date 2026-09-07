package components

import (
	"testing"

	"tuiscale/internal/tailscale"
)

func TestPeerTableModelFilteringAndSorting(t *testing.T) {
	status := &tailscale.Status{
		Self: &tailscale.PeerStatus{
			HostName:     "local-node",
			OS:           "linux",
			TailscaleIPs: []string{"100.64.0.1"},
			Active:       true,
		},
		Peer: map[string]*tailscale.PeerStatus{
			"p1": {
				HostName:     "server-alpha",
				OS:           "linux",
				TailscaleIPs: []string{"100.64.0.10"},
				Active:       true,
			},
			"p2": {
				HostName:     "phone-beta",
				OS:           "android",
				TailscaleIPs: []string{"100.64.0.20"},
				Online:       true,
				Active:       false,
			},
			"p3": {
				HostName:     "desktop-gamma",
				OS:           "windows",
				TailscaleIPs: []string{"100.64.0.30"},
				Online:       false,
				Active:       false,
			},
		},
	}

	model := NewPeerTableModel()

	// 1. Atualizar sem filtro
	model.UpdatePeers(status, "")
	if len(model.peers) != 4 {
		t.Fatalf("Expected 4 peers, got %d", len(model.peers))
	}

	// Active deve vir antes de Offline
	first := model.peers[0]
	last := model.peers[len(model.peers)-1]
	if !first.Active {
		t.Errorf("First peer should be active, got %s", first.HostName)
	}
	if last.Active || last.Online {
		t.Errorf("Last peer should be offline, got %s", last.HostName)
	}

	// 2. Filtro por SO
	model.UpdatePeers(status, "android")
	if len(model.peers) != 1 || model.peers[0].HostName != "phone-beta" {
		t.Errorf("Expected only phone-beta, got %d peers", len(model.peers))
	}

	// 3. Filtro por IP
	model.UpdatePeers(status, "100.64.0.30")
	if len(model.peers) != 1 || model.peers[0].HostName != "desktop-gamma" {
		t.Errorf("Expected desktop-gamma for IP search, got %d peers", len(model.peers))
	}

	// 4. Navegação e limites
	model.UpdatePeers(status, "")
	model.MoveTop()
	if model.cursor != 0 {
		t.Errorf("Expected cursor at 0, got %d", model.cursor)
	}
	model.MoveUp() // Deve permanecer em 0
	if model.cursor != 0 {
		t.Errorf("Expected cursor at 0 after MoveUp at top, got %d", model.cursor)
	}
	model.MoveBottom()
	if model.cursor != 3 {
		t.Errorf("Expected cursor at 3, got %d", model.cursor)
	}
	model.MoveDown() // Deve permanecer em 3
	if model.cursor != 3 {
		t.Errorf("Expected cursor at 3 after MoveDown at bottom, got %d", model.cursor)
	}
}
