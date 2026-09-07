package tailscale

import (
	"math"
	"testing"
	"time"
)

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0 B/s"},
		{-10, "0 B/s"},
		{math.NaN(), "0 B/s"},
		{500, "500 B/s"},
		{1024, "1.0 KB/s"},
		{1536, "1.5 KB/s"},
		{1024 * 1024, "1.00 MB/s"},
		{2.5 * 1024 * 1024, "2.50 MB/s"},
		{1024 * 1024 * 1024, "1.00 GB/s"},
	}

	for _, tt := range tests {
		got := FormatSpeed(tt.input)
		if got != tt.expected {
			t.Errorf("FormatSpeed(%v) = %s; want %s", tt.input, got, tt.expected)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0 B"},
		{-5, "0 B"},
		{800, "800 B"},
		{1024, "1.0 KB"},
		{1024 * 1024 * 5, "5.0 MB"},
		{1024 * 1024 * 1024 * 2, "2.0 GB"},
	}

	for _, tt := range tests {
		got := FormatBytes(tt.input)
		if got != tt.expected {
			t.Errorf("FormatBytes(%v) = %s; want %s", tt.input, got, tt.expected)
		}
	}
}

func TestFormatLatency(t *testing.T) {
	if got := FormatLatency(0); got != "-" {
		t.Errorf("FormatLatency(0) = %s; want '-'", got)
	}
	if got := FormatLatency(50 * time.Millisecond); got != "50 ms" {
		t.Errorf("FormatLatency(50ms) = %s; want '50 ms'", got)
	}
	if got := FormatLatency(500 * time.Microsecond); got != "0.5 ms" {
		t.Errorf("FormatLatency(500µs) = %s; want '0.5 ms'", got)
	}
}

func TestSpeedTracker(t *testing.T) {
	st := NewSpeedTracker()

	status1 := &Status{
		Self: &PeerStatus{
			RxBytes: 1000,
			TxBytes: 2000,
		},
		Peer: map[string]*PeerStatus{
			"p1": {RxBytes: 500, TxBytes: 800},
		},
	}

	st.Update(status1)

	if st.CumulativeRx != 1000 || st.CumulativeTx != 2000 {
		t.Errorf("Cumulative initial failed: got Rx=%d, Tx=%d", st.CumulativeRx, st.CumulativeTx)
	}

	// Simular passagem de tempo
	time.Sleep(100 * time.Millisecond)

	status2 := &Status{
		Self: &PeerStatus{
			RxBytes: 2000, // +1000 bytes
			TxBytes: 4000, // +2000 bytes
		},
		Peer: map[string]*PeerStatus{
			"p1": {RxBytes: 1000, TxBytes: 1600},
		},
	}

	st.Update(status2)

	if st.CurrentRx <= 0 || st.CurrentTx <= 0 {
		t.Errorf("Expected positive speeds, got Rx=%f, Tx=%f", st.CurrentRx, st.CurrentTx)
	}
}
