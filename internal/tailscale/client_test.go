package tailscale

import (
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
