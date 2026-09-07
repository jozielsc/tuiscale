# TUIScale

Interface de Terminal (TUI) moderna, rápida e elegante escrita em **Go** para monitoramento e controle do **Tailscale** no Linux.

Construído com [Bubble Tea](https://github.com/charmbracelet/bubbletea) e [Lip Gloss](https://github.com/charmbracelet/lipgloss), com suporte a compilação estática (`CGO_ENABLED=0`), garantindo compatibilidade universal com qualquer distribuição Linux.

---

## ✨ Recursos

- 📊 **Dashboard Completo de Status**:
  - Estado da conexão (`Conectado`, `Desconectado`, `Login Necessário`).
  - IPs Tailscale (IPv4 e IPv6) e domínio MagicDNS.
  - Nome do nó local, usuário conectado e informações do tailnet.
- ⚡ **Medidor de Velocidade em Tempo Real**:
  - Cálculo instantâneo da taxa de transferência de **Download (Rx)** e **Upload (Tx)** em KB/s e MB/s.
  - Tráfego de dados acumulado total e por dispositivo.
- 🖥️ **Navegador e Inspetor de Peers**:
  - Lista completa de dispositivos da rede com ícones de estado:
    - 🟢 Conexão direta ativa (Direct UDP).
    - 🔵 Conectado via Relay DERP.
    - ⚪ Dispositivo offline.
  - Identificação de Sistema Operacional (🐧 Linux, 🪟 Windows, 🍎 macOS, 🤖 Android, 📱 iOS).
  - Filtro e busca instantânea (`/`) por nome, IP, SO ou usuário.
  - Painel de inspeção aprofundado com detalhes de chave, rotas e endpoints STUN/NAT.
- 🌐 **Diagnóstico de Rede & Netcheck**:
  - Visualização de regiões DERP globais ordenadas por menor latência/RTT.
  - Status de conectividade UDP, IPv4 e IPv6.
  - Detecção de tipo de NAT (estável vs simétrico) e mapeamento de portas (UPnP, PMP, PCP).
  - Verificação de Captive Portal e IP público.
- 🎯 **Ações Rápidas no Terminal**:
  - **Conectar** (`c` -> `tailscale up`) e **Desconectar** (`d` -> `tailscale down`).
  - **Ping interativo** (`p`) para testar latência e caminho de roteamento até o peer selecionado.
  - **Gerenciamento de Nó de Saída / Exit Node** (`e`).
  - **Cópia rápida de IP** (`y`) ou MagicDNS (`Y`) para a área de transferência usando sequências ANSI OSC 52 (compatível nativamente mesmo em sessões SSH e multiplexadores como tmux).
- 🎨 **Tema Dark Mode Nativo**:
  - Visual moderno e refinado com cores e contrastes inspirados no Tailscale e Midnight Dark.

---

## 🐧 Compatibilidade com Distribuições Linux

O executável é gerado como um binário ELF estático autossuficiente (`CGO_ENABLED=0`), funcionando sem alterações em:
- **Debian / Ubuntu / Linux Mint / Pop!_OS**
- **Fedora / RHEL / CentOS / Rocky Linux / AlmaLinux**
- **Arch Linux / Manjaro**
- **Void Linux** (glibc e musl)
- **Alpine Linux**
- **openSUSE / SLES**

---

## ⚙️ Pré-requisitos & Permissões

O **Tailscale** deve estar instalado e o daemon `tailscaled` em execução no sistema.

### Dica de Permissão (Sem necessidade de sudo):
Para permitir que o seu usuário comum controle conexões (`up`, `down`, `set exit-node`) sem precisar digitar `sudo`:
```bash
sudo tailscale set --operator=$USER
```

---

## 🚀 Instalação e Compilação

### Compilação Rápida com `make`:
```bash
# Compilar binário padrão
make build

# Ou compilar binário estático universal (recomendado para empacotamento multi-distro)
make build-static
```

### Instalar no Sistema (`/usr/local/bin`):
```bash
sudo make install
```

### Executar diretamente:
```bash
./tuiscale
```

---

## ⌨️ Atalhos de Teclado

| Tecla | Ação |
| :--- | :--- |
| `↑` / `k`, `↓` / `j` | Navegar pela lista de dispositivos |
| `PgUp` / `PgDn` | Subir ou descer página inteira |
| `g` / `G` | Ir para o topo ou final da lista |
| `Tab` / `Shift+Tab` | Alternar entre as abas (`[1] Peers`, `[2] Rede`, `[3] Métricas`) |
| `1`, `2`, `3` | Ir diretamente para a aba desejada |
| `Enter` | Abrir inspeção detalhada do peer selecionado |
| `/` | Ativar barra de busca / filtro rápido |
| `Esc` | Limpar filtro ativo ou fechar janelas modais |
| `c` | Conectar ao Tailscale (`tailscale up`) |
| `d` | Desconectar do Tailscale (`tailscale down`) |
| `p` | Disparar teste de ping contra o dispositivo selecionado |
| `e` | Abrir seletor de Nó de Saída (Exit Node) |
| `y` | Copiar endereço IP Tailscale para a área de transferência |
| `Y` | Copiar nome de domínio MagicDNS para a área de transferência |
| `n` | Executar diagnóstico de rede (`tailscale netcheck`) |
| `r` | Forçar atualização de status imediata |
| `?` | Exibir tela de ajuda com todos os atalhos |
| `q` / `Ctrl+C` | Sair do aplicativo |

---

## 🛠️ Opções de Linha de Comando (Flags)

```bash
tuiscale [flags]

Flags:
  -refresh duration
        Intervalo de atualização das métricas e status (padrão: 2s)
  -socket string
        Caminho alternativo para o socket do tailscaled (padrão: /var/run/tailscale/tailscaled.sock)
  -v, -version
        Exibe a versão do TUIScale e encerra
```

---

## 📄 Licença

MIT License.
