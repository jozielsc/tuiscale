package clipboard

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
)

// Copy envia o texto para a área de transferência usando OSC 52 e ferramentas do sistema como fallback.
func Copy(text string) error {
	if text == "" {
		return nil
	}

	// 1. Tentar cópia via sequência ANSI OSC 52 no terminal (funciona local e via SSH)
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	osc52Seq := fmt.Sprintf("\x1b]52;c;%s\x07", encoded)
	// Também enviar formato para tmux se estiver rodando dentro do tmux
	if os.Getenv("TMUX") != "" {
		osc52Seq = fmt.Sprintf("\x1bPtmux;\x1b\x1b]52;c;%s\x07\x1b\\", encoded)
	}
	_, _ = os.Stdout.WriteString(osc52Seq)

	// 2. Tentar utilitários nativos de Linux (Wayland wl-copy ou X11 xclip/xsel) silenciosamente
	if _, err := exec.LookPath("wl-copy"); err == nil {
		cmd := exec.Command("wl-copy")
		cmd.Stdin = os.Stdin
		go func() {
			c := exec.Command("wl-copy")
			c.Stdin = os.Stdin
			_ = c.Run()
		}()
	} else if _, err := exec.LookPath("xclip"); err == nil {
		go func() {
			c := exec.Command("xclip", "-selection", "clipboard")
			c.Stdin = os.Stdin
			_ = c.Run()
		}()
	}

	return nil
}
