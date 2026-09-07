package tailscale

import (
	"context"
	"errors"
	"testing"
)

func TestIsDaemonNotRunningError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errors.New("failed to connect to local tailscaled (which appears to be running as tailscaled, pid 615). Got error: Failed to connect to local Tailscale daemon"), true},
		{errors.New("dial unix /var/run/tailscale/tailscaled.sock: connect: no such file or directory"), true},
		{errors.New("connect: connection refused"), true},
		{errors.New("Tailscale daemon is not running"), true},
		{errors.New("is tailscaled running?"), true},
		{errors.New("random unrelated error"), false},
	}

	for _, tt := range tests {
		got := IsDaemonNotRunningError(tt.err)
		if got != tt.expected {
			t.Errorf("IsDaemonNotRunningError(%v) = %v; want %v", tt.err, got, tt.expected)
		}
	}
}

func TestBuildCommand(t *testing.T) {
	// Teste sem socket customizado
	client := NewClient("")
	cmd := client.buildCommand(context.Background(), "status", "--json")
	if cmd.Path == "" {
		t.Error("expected command path to be set")
	}

	// Teste com socket customizado
	clientSocket := NewClient("/custom/socket/path")
	cmdSocket := clientSocket.buildCommand(context.Background(), "status")

	// Verificar se o socket foi adicionado aos argumentos
	found := false
	for _, arg := range cmdSocket.Args {
		if arg == "--socket=/custom/socket/path" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected socket flag in command args")
	}
}

func TestClientImplementsInterface(t *testing.T) {
	// Teste de compilação para garantir que Client implementa TailscaleClient
	var _ TailscaleClient = NewClient("")
}
