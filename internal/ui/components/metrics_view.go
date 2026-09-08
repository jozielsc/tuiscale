package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/tailscale"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// RenderMetricsView desenha o painel com estatísticas de tráfego, contadores e distribuição por dispositivo.
func RenderMetricsView(status *tailscale.Status, speedTracker *tailscale.SpeedTracker, width, height int) string {
	if status == nil {
		return theme.Panel.Width(width).Height(height).Render("Carregando métricas do Tailscale...")
	}

	// 1. Contadores gerais
	totalPeers := len(status.Peer)
	var onlinePeers, directPeers, relayPeers, offlinePeers int
	for _, p := range status.Peer {
		if p.Active {
			directPeers++
			onlinePeers++
		} else if p.Online {
			relayPeers++
			onlinePeers++
		} else {
			offlinePeers++
		}
	}

	tunStr := "Desabilitado"
	if status.TUN {
		tunStr = "Habilitado (TUN Device)"
	}

	keyStatus := "Válida"
	if !status.HaveNodeKey {
		keyStatus = "Ausente / Expirada"
	}

	// Labels e formato adaptados à largura disponível
	var labelFmt string
	if width >= 80 {
		labelFmt = "%-22s"
	} else {
		labelFmt = "%-16s"
		// Versão compacta de alguns valores para telas estreitas
		tunStr = "Habilitado"
		if !status.TUN {
			tunStr = "Desabilitado"
		}
	}

	secSummary := fmt.Sprintf(labelFmt+" %d nós na tailnet\n"+labelFmt+" %s\n"+labelFmt+" %s\n"+labelFmt+" %s",
		"Total Dispositivos:", totalPeers,
		"Status dos Nós:", fmt.Sprintf("🟢 %d Dir  🔵 %d Relay  ⚪ %d Off", directPeers, relayPeers, offlinePeers),
		"Rede (TUN):", tunStr,
		"Chave:", keyStatus,
	)

	// 2. Tráfego Atual e Acumulado
	var curRx, curTx float64
	var cumRx, cumTx int64
	if speedTracker != nil {
		curRx = speedTracker.CurrentRx
		curTx = speedTracker.CurrentTx
		cumRx = speedTracker.CumulativeRx
		cumTx = speedTracker.CumulativeTx
	}

	secTraffic := fmt.Sprintf(labelFmt+" %s\n"+labelFmt+" %s\n"+labelFmt+" %s\n"+labelFmt+" %s",
		"Download (Rx):", theme.SpeedRxStyle.Render(tailscale.FormatSpeed(curRx)),
		"Upload (Tx):", theme.SpeedTxStyle.Render(tailscale.FormatSpeed(curTx)),
		"Total Download:", tailscale.FormatBytes(cumRx),
		"Total Upload:", tailscale.FormatBytes(cumTx),
	)

	// Layout dos cards: lado a lado para width>=80, empilhado para width<80
	var cardBlock string
	if width >= 80 {
		cardBlock = lipgloss.JoinHorizontal(lipgloss.Top,
			theme.Panel.Width((width/2)-2).Render(theme.TitleStyle.Render("RESUMO DA REDE")+"\n\n"+secSummary),
			theme.Panel.Width((width/2)-2).Render(theme.TitleStyle.Render("TRÁFEGO DE DADOS")+"\n\n"+secTraffic),
		)
	} else {
		cardBlock = lipgloss.JoinVertical(lipgloss.Left,
			theme.Panel.Width(width-4).Render(theme.TitleStyle.Render("RESUMO DA REDE")+"\n\n"+secSummary),
			theme.Panel.Width(width-4).Render(theme.TitleStyle.Render("TRÁFEGO DE DADOS")+"\n\n"+secTraffic),
		)
	}

	// 3. Top Dispositivos por Consumo de Tráfego (Rx + Tx)
	type peerTraffic struct {
		name       string
		totalBytes int64
		rxBytes    int64
		txBytes    int64
	}

	var trafficList []peerTraffic
	for _, p := range status.Peer {
		tot := p.RxBytes + p.TxBytes
		if tot > 0 {
			trafficList = append(trafficList, peerTraffic{
				name:       p.HostName,
				totalBytes: tot,
				rxBytes:    p.RxBytes,
				txBytes:    p.TxBytes,
			})
		}
	}

	sort.Slice(trafficList, func(i, j int) bool {
		return trafficList[i].totalBytes > trafficList[j].totalBytes
	})

	var maxBytes int64 = 1
	if len(trafficList) > 0 && trafficList[0].totalBytes > 0 {
		maxBytes = trafficList[0].totalBytes
	}

	var rankingRows []string
	// Layout do ranking adaptado à largura: em telas estreitas omite a coluna
	// de proporção (barra) e reduz larguras das colunas de nome/tráfego.
	narrowRanking := width < 80
	var rankingHeader string
	if narrowRanking {
		rankingHeader = fmt.Sprintf("  %-18s %-12s %-12s", "DISPOSITIVO", "RX", "TX")
	} else {
		rankingHeader = fmt.Sprintf("  %-22s %-16s %-16s %-24s", "DISPOSITIVO", "DOWNLOAD (RX)", "UPLOAD (TX)", "PROPORÇÃO")
	}
	rankingRows = append(rankingRows, theme.TableHeader.Render(rankingHeader))

	maxRows := height - 14
	if maxRows < 4 {
		maxRows = 4
	}
	if len(trafficList) > maxRows {
		trafficList = trafficList[:maxRows]
	}

	for _, pt := range trafficList {
		if narrowRanking {
			// Sem barra de proporção; nome truncado
			row := fmt.Sprintf("  %-18s %-12s %-12s",
				truncate(pt.name, 18),
				tailscale.FormatBytes(pt.rxBytes),
				tailscale.FormatBytes(pt.txBytes),
			)
			rankingRows = append(rankingRows, row)
			continue
		}
		barLen := 18
		ratio := float64(pt.totalBytes) / float64(maxBytes)
		filled := int(ratio * float64(barLen))
		if filled < 1 {
			filled = 1
		}
		empty := barLen - filled
		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
		barColored := lipgloss.NewStyle().Foreground(theme.ColorHighlight).Render("[" + bar + "]")

		row := fmt.Sprintf("  %-22s %-16s %-16s %s",
			pt.name,
			tailscale.FormatBytes(pt.rxBytes),
			tailscale.FormatBytes(pt.txBytes),
			barColored,
		)
		rankingRows = append(rankingRows, row)
	}

	rankingBlock := ""
	if len(trafficList) > 0 {
		rankingBlock = fmt.Sprintf("\n%s\n%s",
			theme.SubTitleStyle.Render("Top Dispositivos por Consumo de Tráfego:"),
			strings.Join(rankingRows, "\n"),
		)
	}

	content := fmt.Sprintf("%s\n\n%s%s",
		theme.TitleStyle.Render("MÉTRICAS E ESTATÍSTICAS DO TAILSCALE"),
		cardBlock,
		rankingBlock,
	)

	return theme.Panel.Width(width).Height(height).Render(content)
}
