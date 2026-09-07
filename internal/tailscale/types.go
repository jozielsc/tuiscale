package tailscale

import "time"

// Status representa a saída JSON de `tailscale status --json`.
type Status struct {
	Version        string                   `json:"Version"`
	TUN            bool                     `json:"TUN"`
	BackendState   string                   `json:"BackendState"`
	HaveNodeKey    bool                     `json:"HaveNodeKey"`
	AuthURL        string                   `json:"AuthURL"`
	TailscaleIPs   []string                 `json:"TailscaleIPs"`
	Self           *PeerStatus              `json:"Self"`
	Health         []string                 `json:"Health"`
	MagicDNSSuffix string                   `json:"MagicDNSSuffix"`
	CurrentTailnet *Tailnet                 `json:"CurrentTailnet"`
	Peer           map[string]*PeerStatus   `json:"Peer"`
	User           map[string]TailscaleUser `json:"User"`
}

// Tailnet contém detalhes sobre a rede tailnet conectada.
type Tailnet struct {
	Name            string `json:"Name"`
	MagicDNSSuffix  string `json:"MagicDNSSuffix"`
	MagicDNSEnabled bool   `json:"MagicDNSEnabled"`
}

// TailscaleUser representa um usuário na tailnet.
type TailscaleUser struct {
	ID          int64  `json:"ID"`
	LoginName   string `json:"LoginName"`
	DisplayName string `json:"DisplayName"`
}

// PeerStatus representa as propriedades de um nó (local ou remoto).
type PeerStatus struct {
	ID             string   `json:"ID"`
	NodeID         int64    `json:"NodeID"`
	PublicKey      string   `json:"PublicKey"`
	HostName       string   `json:"HostName"`
	DNSName        string   `json:"DNSName"`
	OS             string   `json:"OS"`
	UserID         int64    `json:"UserID"`
	TailscaleIPs   []string `json:"TailscaleIPs"`
	AllowedIPs     []string `json:"AllowedIPs"`
	Addrs          []string `json:"Addrs"`
	CurAddr        string   `json:"CurAddr"`
	Relay          string   `json:"Relay"`
	RxBytes        int64    `json:"RxBytes"`
	TxBytes        int64    `json:"TxBytes"`
	Created        string   `json:"Created"`
	LastWrite      string   `json:"LastWrite"`
	LastSeen       string   `json:"LastSeen"`
	LastHandshake  string   `json:"LastHandshake"`
	Online         bool     `json:"Online"`
	ExitNode       bool     `json:"ExitNode"`
	ExitNodeOption bool     `json:"ExitNodeOption"`
	Active         bool     `json:"Active"`
	KeyExpiry      string   `json:"KeyExpiry"`

	// Campos calculados em tempo de execução
	UserLogin  string  `json:"-"`
	SpeedRx    float64 `json:"-"` // Bytes por segundo
	SpeedTx    float64 `json:"-"` // Bytes por segundo
	LastPingMs float64 `json:"-"` // Última latência medida via ping
}

// NetcheckReport representa a resposta JSON de `tailscale netcheck --format=json`.
type NetcheckReport struct {
	Now                   string             `json:"Now"`
	UDP                   bool               `json:"UDP"`
	IPv4                  bool               `json:"IPv4"`
	IPv6                  bool               `json:"IPv6"`
	IPv4CanSend           bool               `json:"IPv4CanSend"`
	IPv6CanSend           bool               `json:"IPv6CanSend"`
	OSHasIPv6             bool               `json:"OSHasIPv6"`
	ICMPv4                bool               `json:"ICMPv4"`
	MappingVariesByDestIP bool               `json:"MappingVariesByDestIP"`
	UPnP                  bool               `json:"UPnP"`
	PMP                   bool               `json:"PMP"`
	PCP                   bool               `json:"PCP"`
	PreferredDERP         int                `json:"PreferredDERP"`
	RegionLatency         map[string]float64 `json:"RegionLatency"`   // em nanossegundos
	RegionV4Latency       map[string]float64 `json:"RegionV4Latency"` // em nanossegundos
	GlobalV4              string             `json:"GlobalV4"`
	GlobalV6              string             `json:"GlobalV6"`
	CaptivePortal         bool               `json:"CaptivePortal"`
}

// SpeedSnapshot armazena as taxas de tráfego atuais (Rx/Tx) agregadas e por peer.
type SpeedSnapshot struct {
	Timestamp   time.Time
	TotalRxRate float64 // Bytes por segundo recebidos
	TotalTxRate float64 // Bytes por segundo enviados
	TotalRx     int64
	TotalTx     int64
}

// DERPRegionInfo descreve a localização e nome comum de um DERP ID conhecido.
type DERPRegionInfo struct {
	ID       int
	Code     string
	City     string
	Country  string
	Latency  time.Duration
	Selected bool
}
