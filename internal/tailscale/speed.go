package tailscale

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// SpeedTracker calcula taxas de transferência (download e upload) ao longo do tempo.
// Ele mantém snapshots dos contadores anteriores para derivar a taxa instantânea entre chamadas a Update.
// É seguro para uso concorrente via mutex interno.
type SpeedTracker struct {
	mu           sync.Mutex
	lastCheck    time.Time
	prevSelfRx   int64
	prevSelfTx   int64
	prevPeerRx   map[string]int64
	prevPeerTx   map[string]int64
	CurrentRx    float64 // B/s
	CurrentTx    float64 // B/s
	CumulativeRx int64
	CumulativeTx int64
}

// NewSpeedTracker cria uma nova instância de SpeedTracker pronta para uso.
func NewSpeedTracker() *SpeedTracker {
	return &SpeedTracker{
		prevPeerRx: make(map[string]int64),
		prevPeerTx: make(map[string]int64),
	}
}

// Update calcula as taxas instantâneas de Rx/Tx com base no novo Status retornado pelo Tailscale.
// Na primeira chamada apenas inicializa os contadores de baseline.
// Protege os campos do SpeedTracker com mutex; os campos SpeedRx/SpeedTx dos peers
// são escritos diretamente pois pertencem ao status recém-criado do ciclo de atualização.
func (st *SpeedTracker) Update(status *Status) {
	if status == nil {
		return
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	now := time.Now()

	// Se não tiver Self ou primeira execução, apenas inicializamos os contadores
	if st.lastCheck.IsZero() {
		st.lastCheck = now
		if status.Self != nil {
			st.prevSelfRx = status.Self.RxBytes
			st.prevSelfTx = status.Self.TxBytes
			st.CumulativeRx = status.Self.RxBytes
			st.CumulativeTx = status.Self.TxBytes
		}
		for id, p := range status.Peer {
			if p != nil {
				st.prevPeerRx[id] = p.RxBytes
				st.prevPeerTx[id] = p.TxBytes
			}
		}
		return
	}

	dt := now.Sub(st.lastCheck).Seconds()
	if dt <= 0 {
		return
	}
	st.lastCheck = now

	// 1. Cálculo para o nó local (Self)
	if status.Self != nil {
		deltaRx := status.Self.RxBytes - st.prevSelfRx
		deltaTx := status.Self.TxBytes - st.prevSelfTx

		// Em caso de reinicialização do daemon, o delta pode ser negativo
		if deltaRx < 0 {
			deltaRx = 0
		}
		if deltaTx < 0 {
			deltaTx = 0
		}

		st.CurrentRx = float64(deltaRx) / dt
		st.CurrentTx = float64(deltaTx) / dt
		st.prevSelfRx = status.Self.RxBytes
		st.prevSelfTx = status.Self.TxBytes
		st.CumulativeRx = status.Self.RxBytes
		st.CumulativeTx = status.Self.TxBytes
	}

	// 2. Cálculo por Peer individual
	for id, p := range status.Peer {
		if p == nil {
			continue
		}
		prevRx := st.prevPeerRx[id]
		prevTx := st.prevPeerTx[id]

		deltaRx := p.RxBytes - prevRx
		deltaTx := p.TxBytes - prevTx
		if deltaRx < 0 {
			deltaRx = 0
		}
		if deltaTx < 0 {
			deltaTx = 0
		}

		p.SpeedRx = float64(deltaRx) / dt
		p.SpeedTx = float64(deltaTx) / dt

		st.prevPeerRx[id] = p.RxBytes
		st.prevPeerTx[id] = p.TxBytes
	}
}

// FormatSpeed formata bytes por segundo em formato legível (e.g. 1.2 MB/s).
func FormatSpeed(bytesPerSec float64) string {
	if bytesPerSec <= 0 || math.IsNaN(bytesPerSec) || math.IsInf(bytesPerSec, 0) {
		return "0 B/s"
	}
	const (
		kb = 1024.0
		mb = 1024.0 * kb
		gb = 1024.0 * mb
	)

	switch {
	case bytesPerSec >= gb:
		return fmt.Sprintf("%.2f GB/s", bytesPerSec/gb)
	case bytesPerSec >= mb:
		return fmt.Sprintf("%.2f MB/s", bytesPerSec/mb)
	case bytesPerSec >= kb:
		return fmt.Sprintf("%.1f KB/s", bytesPerSec/kb)
	default:
		return fmt.Sprintf("%.0f B/s", bytesPerSec)
	}
}

// FormatBytes formata bytes brutos em formato legível (e.g. 23.4 MB).
func FormatBytes(b int64) string {
	if b <= 0 {
		return "0 B"
	}
	const unit = 1024.0
	if b < int64(unit) {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / int64(unit); n >= int64(unit); n /= int64(unit) {
		div *= int64(unit)
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}
	return fmt.Sprintf("%.1f %s", float64(b)/float64(div), units[exp])
}

// FormatLatency formata duração em milissegundos legíveis.
func FormatLatency(d time.Duration) string {
	if d <= 0 {
		return "-"
	}
	ms := float64(d.Microseconds()) / 1000.0
	if ms < 1.0 {
		return fmt.Sprintf("%.1f ms", ms)
	}
	return fmt.Sprintf("%.0f ms", ms)
}
