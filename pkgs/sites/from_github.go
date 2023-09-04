package sites

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/moqsien/goutils/pkgs/crypt"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gscraper/pkgs/conf"
	"github.com/moqsien/vpnparser/pkgs/outbound"
)

var (
	VPN_MAP = map[string]struct{}{}
)

type VPNFromGithub struct {
	Conf    *conf.Config
	VPNList *outbound.Result `json:"vpn_list"`
}

func NewVPNFromGithub(cnf *conf.Config, vpnList *outbound.Result) (v *VPNFromGithub) {
	v = &VPNFromGithub{
		Conf:    cnf,
		VPNList: vpnList,
	}
	v.initiate()
	return v
}

func (that *VPNFromGithub) initiate() {
	if that.Conf.NeoboxKey == "" {
		that.Conf.ResetNeoboxKey()
		that.Conf.ShowGithubVPNSubscriber()
	}
}

func (that *VPNFromGithub) decodeBase64(content []byte) (result string) {
	rawStr := string(content)
	if strings.Contains(rawStr, "://") {
		result = rawStr
	} else {
		if strings.Contains(rawStr, "\n") {
			rList := strings.Split(rawStr, "\n")
			for _, r := range rList {
				rr := crypt.DecodeBase64(r)
				if result == "" && rr != "" {
					result = rr
				}
				if result != "" && rr != "" {
					result = result + "\n" + rr
				}
			}
		} else {
			result = crypt.DecodeBase64(rawStr)
		}
	}
	return
}

func (that *VPNFromGithub) Parse(content []byte) {
	res := that.decodeBase64(content)
	if res == "" {
		tui.PrintError("base64 decoding failed")
		rawStr := string(content)
		if len(rawStr) > 5 {
			fmt.Println(rawStr[len(rawStr)-4:])
		} else {
			fmt.Println(rawStr)
		}
		return
	}

	vList := strings.Split(res, "\n")
	for _, v := range vList {
		v = strings.TrimSpace(v)
		v = strings.TrimRight(v, "\r")
		v = ParseRawUri(v)
		if strings.HasPrefix(v, "vmess://") {
			_, ok := VPN_MAP[v]
			if s := that.wrapItem(v); !ok && s != nil {
				that.VPNList.Vmess = append(that.VPNList.Vmess, s)
				VPN_MAP[v] = struct{}{}
			}
		} else if strings.HasPrefix(v, "vless://") {
			_, ok := VPN_MAP[v]
			if s := that.wrapItem(v); !ok && s != nil {
				that.VPNList.Vless = append(that.VPNList.Vless, s)
				VPN_MAP[v] = struct{}{}
			}
		} else if strings.HasPrefix(v, "ss://") {
			_, ok := VPN_MAP[v]
			if s := that.wrapItem(v); !ok && s != nil {
				that.VPNList.ShadowSocks = append(that.VPNList.ShadowSocks, s)
				VPN_MAP[v] = struct{}{}
			}
		} else if strings.HasPrefix(v, "ssr://") {
			_, ok := VPN_MAP[v]
			if s := that.wrapItem(v); !ok && s != nil {
				that.VPNList.ShadowSocksR = append(that.VPNList.ShadowSocksR, s)
				VPN_MAP[v] = struct{}{}
			}
		} else if strings.HasPrefix(v, "trojan://") {
			_, ok := VPN_MAP[v]
			if s := that.wrapItem(v); !ok && s != nil {
				that.VPNList.Trojan = append(that.VPNList.Trojan, s)
				VPN_MAP[v] = struct{}{}
			}

		}
	}
}

func (that *VPNFromGithub) wrapItem(rawUri string) *outbound.ProxyItem {
	item := outbound.NewItem(rawUri)
	if strings.HasPrefix(item.Address, "127.0.") {
		return nil
	}
	item.GetOutbound()
	return item
}

func (that *VPNFromGithub) getUrl(sUrl string) (pUrl string) {
	if strings.Contains(sUrl, "raw.githubusercontent.com") {
		pUrl = that.Conf.GithubSpeedupUrl + sUrl
	} else if strings.Contains(sUrl, "{year}") || strings.Contains(sUrl, "{month}") || strings.Contains(sUrl, "{day}") {
		now := time.Now()
		sUrl = strings.ReplaceAll(sUrl, "{year}", now.Format("2006"))
		sUrl = strings.ReplaceAll(sUrl, "{month}", now.Format("01"))
		pUrl = strings.ReplaceAll(sUrl, "{day}", now.Format("02"))
	} else {
		pUrl = sUrl
	}
	return
}

func (that *VPNFromGithub) Run() {
	for _, sUrl := range that.Conf.GithubVPNSubscriber {
		client := &http.Client{Timeout: 30 * time.Second}
		pUrl := that.getUrl(sUrl)
		req, err := http.NewRequest("GET", pUrl, nil)
		if err != nil {
			tui.PrintError(err)
			continue
		}
		tui.PrintInfo(pUrl)
		rep, err := client.Do(req)
		if err != nil {
			tui.PrintError(err)
			continue
		}
		defer rep.Body.Close()

		data, err := io.ReadAll(rep.Body)
		// fmt.Println(rep.ContentLength, len(data), rep.StatusCode)
		if err != nil || rep.StatusCode != 200 {
			tui.PrintError(fmt.Sprintf("StatusCode: %v", rep.StatusCode), err)
			continue
		}

		if len(data) > 0 {
			that.Parse(data)
		}
	}
}
