package sites

type Result struct {
	Vmess        []string `json:"vmess"`
	Vless        []string `json:"vless"`
	ShadowSocks  []string `json:"shadowsocks"`
	ShadowSocksR []string `json:"shadowsocksR"`
	Trojan       []string `json:"trojan"`
	UpdateAt     string   `json:"update_time"`
}
