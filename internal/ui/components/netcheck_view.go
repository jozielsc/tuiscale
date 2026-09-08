package components

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/jozielsc/tuiscale/internal/tailscale"
	"github.com/jozielsc/tuiscale/internal/ui/theme"
)

// RenderNetcheckView desenha o painel de diagnóstico completo do Tailscale (DERP, NAT, IPv4/IPv6).
func RenderNetcheckView(report *tailscale.NetcheckReport, isRunning bool, width, height int) string {
	if isRunning {
		return theme.Panel.Width(width).Height(height).Render(
			lipgloss.NewStyle().
				Width(width).
				Height(height).
				Align(lipgloss.Center, lipgloss.Center).
				Render("Executando diagnóstico de rede (netcheck)... Aguarde alguns instantes."),
		)
	}

	if report == nil {
		return theme.Panel.Width(width).Height(height).Render(
			lipgloss.NewStyle().
				Width(width).
				Height(height).
				Align(lipgloss.Center, lipgloss.Center).
				Render("Nenhum relatório de netcheck disponível ainda.\nPressione 'n' para disparar o diagnóstico de rede."),
		)
	}

	// 1. Bloco de Diagnóstico Geral
	boolBadge := func(val bool, yesText, noText string) string {
		if val {
			return theme.BadgeDirect.Render("✓ " + yesText)
		}
		return theme.BadgeOffline.Render("✗ " + noText)
	}

	natType := "Fácil (Mapeamento estável)"
	if report.MappingVariesByDestIP {
		natType = theme.BadgeRelay.Render("Difícil / Simétrico (Varia por destino)")
	}

	prefCode, prefCity, prefCountry := tailscale.GetDERPLocation(report.PreferredDERP)
	prefDERPStr := fmt.Sprintf("#%d [%s] %s, %s", report.PreferredDERP, prefCode, prefCity, prefCountry)

	col1 := fmt.Sprintf("%-20s %s\n%-20s %s\n%-20s %s\n%-20s %s",
		"Conectividade UDP:", boolBadge(report.UDP, "Habilitado", "Bloqueado"),
		"Suporte IPv4:", boolBadge(report.IPv4, "Disponível", "Indisponível"),
		"Suporte IPv6:", boolBadge(report.IPv6, "Disponível", "Indisponível"),
		"Captive Portal:", boolBadge(!report.CaptivePortal, "Nenhum", "Detectado!"),
	)

	col2 := fmt.Sprintf("%-20s %s\n%-20s %s\n%-20s %s\n%-20s %s",
		"IP Público (IPv4):", report.GlobalV4,
		"Tipo de NAT:", natType,
		"Port Mapping:", fmt.Sprintf("UPnP:%t PMP:%t PCP:%t", report.UPnP, report.PMP, report.PCP),
		"DERP Preferido:", lipgloss.NewStyle().Bold(true).Foreground(theme.ColorSuccess).Render(prefDERPStr),
	)

	diagBlock := lipgloss.JoinHorizontal(lipgloss.Top, col1, "     ", col2)

	// 2. Tabela de Latência de Regiões DERP (ordenadas da menor para a maior latência)
	type derpItem struct {
		id      int
		latency time.Duration
		isPref  bool
		city    string
		country string
		code    string
	}

	var derpList []derpItem
	for regStr, latNs := range report.RegionLatency {
		id, _ := strconv.Atoi(regStr)
		lat := time.Duration(latNs)
		code, city, country := tailscale.GetDERPLocation(id)
		derpList = append(derpList, derpItem{
			id:      id,
			latency: lat,
			isPref:  (id == report.PreferredDERP),
			city:    city,
			country: country,
			code:    code,
		})
	}

	sort.Slice(derpList, func(i, j int) bool {
		return derpList[i].latency < derpList[j].latency
	})

	renderLatencyBar := func(d time.Duration) string {
		ms := float64(d.Milliseconds())
		// Normalizar barra em escala de até 400ms
		const barMax = 16
		filled := int((ms / 400.0) * barMax)
		if filled > barMax {
			filled = barMax
		}
		if filled < 1 {
			filled = 1
		}
		empty := barMax - filled

		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
		var barColor lipgloss.Color
		switch {
		case ms < 80:
			barColor = theme.ColorSuccess
		case ms < 180:
			barColor = theme.ColorWarning
		default:
			barColor = theme.ColorDanger
		}
		return lipgloss.NewStyle().Foreground(barColor).Render(fmt.Sprintf("[%s] %3.0f ms", bar, ms))
	}

	var derpRows []string
	derpHeader := fmt.Sprintf("  %-4s %-6s %-16s %-14s %-22s", "ID", "CÓD", "CIDADE", "PAÍS", "LATÊNCIA / RTT")
	derpRows = append(derpRows, theme.TableHeader.Render(derpHeader))

	// Exibir até caber no espaço disponível
	maxRows := height - 12
	if maxRows < 5 {
		maxRows = 5
	}
	if len(derpList) > maxRows {
		derpList = derpList[:maxRows]
	}

	for _, item := range derpList {
		prefMarker := "  "
		if item.isPref {
			prefMarker = "★ "
		}
		row := fmt.Sprintf("%s%-4d %-6s %-16s %-14s %s",
			prefMarker,
			item.id,
			item.code,
			item.city,
			item.country,
			renderLatencyBar(item.latency),
		)
		if item.isPref {
			derpRows = append(derpRows, lipgloss.NewStyle().Bold(true).Foreground(theme.ColorAccent).Render(row))
		} else {
			derpRows = append(derpRows, row)
		}
	}

	content := fmt.Sprintf("%s\n\n%s\n\n%s\n%s",
		theme.TitleStyle.Render("DIAGNÓSTICO DE REDE (TAILSCALE NETCHECK)"),
		diagBlock,
		theme.SubTitleStyle.Render("Latência Global de Servidores DERP (Pressione 'n' para atualizar):"),
		strings.Join(derpRows, "\n"),
	)

	return theme.Panel.Width(width).Height(height).Render(content)
}
