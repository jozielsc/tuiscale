package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"tuiscale/internal/tailscale"
	"tuiscale/internal/ui"
)

var (
	version = "1.0.0"
	commit  = "none"
	date    = "unknown"
)

func main() {
	refreshFlag := flag.Duration("refresh", 2*time.Second, "Intervalo de atualização das métricas e status (ex: 2s, 3s)")
	socketFlag := flag.String("socket", "", "Caminho do socket do tailscaled (padrão /var/run/tailscale/tailscaled.sock)")
	versionFlag := flag.Bool("version", false, "Exibe a versão do TUIScale e sai")
	flag.BoolVar(versionFlag, "v", false, "Exibe a versão do TUIScale e sai (abreviação)")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("TUIScale v%s (commit: %s, build date: %s)\n", version, commit, date)
		return
	}

	// Verificar se o binário do Tailscale está presente no PATH
	if _, err := exec.LookPath("tailscale"); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: O comando 'tailscale' não foi encontrado no seu PATH.\n")
		fmt.Fprintf(os.Stderr, "Certifique-se de que o Tailscale está instalado no sistema.\n")
		fmt.Fprintf(os.Stderr, "Instalação rápida no Linux: curl -fsSL https://tailscale.com/install.sh | sh\n")
		os.Exit(1)
	}

	client := tailscale.NewClient(*socketFlag)
	app := ui.NewAppModel(client, *refreshFlag)

	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro na execução do TUIScale: %v\n", err)
		os.Exit(1)
	}
}
