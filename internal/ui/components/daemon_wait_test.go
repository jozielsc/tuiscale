package components

import (
	"strings"
	"testing"
)

func TestRenderDaemonWaitView(t *testing.T) {
	// 1. Initial connecting state (attempts = 0, no error)
	viewInit := RenderDaemonWaitView(80, 24, "", 0)
	if !strings.Contains(viewInit, "CONECTANDO AO DAEMON") {
		t.Errorf("expected initial state to show 'CONECTANDO AO DAEMON'")
	}

	// 2. Waiting with error and attempts
	errMsg := "dial unix /var/run/tailscale/tailscaled.sock: connect: no such file or directory"
	viewErr := RenderDaemonWaitView(80, 24, errMsg, 2)
	if !strings.Contains(viewErr, "SERVIÇO TAILSCALE INATIVO") {
		t.Errorf("expected 'SERVIÇO TAILSCALE INATIVO'")
	}
	if !strings.Contains(viewErr, "sudo systemctl start tailscaled") {
		t.Errorf("expected suggested command systemctl start")
	}
	if !strings.Contains(viewErr, "sudo rc-service tailscaled start") {
		t.Errorf("expected suggested command rc-service start")
	}
	if !strings.Contains(viewErr, "sudo tailscaled") {
		t.Errorf("expected suggested command sudo tailscaled")
	}
	if !strings.Contains(viewErr, "tentativa 2") {
		t.Errorf("expected attempt 2 mentioned")
	}
	if !strings.Contains(viewErr, "[q] / [Ctrl+C]") {
		t.Errorf("expected quit shortcuts")
	}
	if !strings.Contains(viewErr, "[r]") {
		t.Errorf("expected retry shortcut")
	}
}
