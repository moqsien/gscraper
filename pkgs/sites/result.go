package sites

type Result struct {
	Vmess        []string `json:"vmess"`
	Vless        []string `json:"vless"`
	ShadowSocks  []string `json:"shadowsocks"`
	ShadowSocksR []string `json:"shadowsocksR"`
	Trojan       []string `json:"trojan"`
	UpdateAt     string   `json:"update_time"`
	VmessTotal   int      `json:"vmess_total"`
	VlessTotal   int      `json:"vless_total"`
	TrojanTotal  int      `json:"trojan_total"`
	SSTotal      int      `json:"ss_total"`
	SSRTotal     int      `json:"ssr_total"`
}

var (
	VPN_MAP = map[string]struct{}{}
)
