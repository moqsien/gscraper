package sites

import (
	"encoding/json"
	"strings"

	"github.com/moqsien/vpnparser/pkgs/outbound"
	"github.com/moqsien/vpnparser/pkgs/parser"
	"github.com/moqsien/vpnparser/pkgs/utils"
)

type Result struct {
	Vmess        []string `json:"vmess"`
	Vless        []string `json:"vless"`
	ShadowSocks  []string `json:"shadowsocks"`
	ShadowSocksR []string `json:"shadowsocksR"`
	Trojan       []string `json:"trojan"`
	UpdateAt     string   `json:"update_time"`
}

var (
	VPN_MAP = map[string]struct{}{}
)

type Item struct {
	Address      string `json:"address"`
	Port         int    `json:"port"`
	RawUri       string `json:"raw_uri"`
	Outbound     string `json:"outbound"`
	OutboundType string `json:"outbound_type"`
}

func NewItem(rawUri string) *Item {
	return &Item{RawUri: rawUri}
}

func (that *Item) String() string {
	scheme := utils.ParseScheme(that.RawUri)
	if scheme == parser.SchemeSSR {
		that.OutboundType = string(outbound.SingBox)
		ob := outbound.GetOutbound(outbound.SingBox, that.RawUri)
		if ob == nil {
			return ""
		}
		ob.Parse(that.RawUri)
		that.OutboundType = ob.GetOutboundStr()
		that.Address = ob.Addr()
		that.Port = ob.Port()
	} else if scheme == parser.SchemeSS && strings.Contains(that.RawUri, "plugin=") {
		that.OutboundType = string(outbound.SingBox)
		ob := outbound.GetOutbound(outbound.SingBox, that.RawUri)
		if ob == nil {
			return ""
		}
		ob.Parse(that.RawUri)
		that.OutboundType = ob.GetOutboundStr()
		that.Address = ob.Addr()
		that.Port = ob.Port()
	} else {
		that.OutboundType = string(outbound.XrayCore)
		ob := outbound.GetOutbound(outbound.XrayCore, that.RawUri)
		if ob == nil {
			return ""
		}
		ob.Parse(that.RawUri)
		that.OutboundType = ob.GetOutboundStr()
		that.Address = ob.Addr()
		that.Port = ob.Port()
	}
	if r, err := json.Marshal(that); err == nil {
		return string(r)
	}
	return ""
}
