package sites

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/moqsien/goutils/pkgs/crypt"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/gscraper/pkgs/conf"
)

type VPNFromGithub struct {
	Conf    *conf.Config
	VPNList *Result `json:"vpn_list"`
}

func NewVPNFromGithub(cnf *conf.Config, vpnList *Result) (v *VPNFromGithub) {
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

func (that *VPNFromGithub) Parse(content []byte) {
	res := crypt.DecodeBase64(string(content))
	if res != "" {
		vList := strings.Split(res, "\n")
		for _, v := range vList {
			if strings.HasPrefix(v, "vmess://") {
				that.VPNList.Vmess = append(that.VPNList.Vmess, v)
			} else if strings.HasPrefix(v, "vless://") {
				that.VPNList.Vless = append(that.VPNList.Vless, v)
			} else if strings.HasPrefix(v, "ss://") {
				that.VPNList.ShadowSocks = append(that.VPNList.ShadowSocks, v)
			} else if strings.HasPrefix(v, "ssr://") {
				that.VPNList.ShadowSocksR = append(that.VPNList.ShadowSocksR, v)
			} else if strings.HasPrefix(v, "trojan://") {
				that.VPNList.Trojan = append(that.VPNList.Trojan, v)
			}
		}
	}
}

func (that *VPNFromGithub) Run() {
	for _, sUrl := range that.Conf.GithubVPNSubscriber {
		client := &http.Client{Timeout: 30 * time.Second}
		pUrl := that.Conf.GithubSpeedupUrl + sUrl
		req, err := http.NewRequest("GET", pUrl, nil)
		if err != nil {
			tui.PrintError(err)
			return
		}
		tui.PrintInfo(pUrl)
		rep, err := client.Do(req)
		if err != nil {
			tui.PrintError(err)
			return
		}
		defer rep.Body.Close()

		data, err := io.ReadAll(rep.Body)
		if err != nil {
			tui.PrintError(err)
			return
		}

		if len(data) > 0 {
			that.Parse(data)
		}
	}
}
