package sites

import "github.com/moqsien/gscraper/pkgs/conf"

type IVpnSite interface {
	Parse([]byte)
	Run()
}

type Sites struct {
	conf     *conf.Config
	siteList []IVpnSite
	VPNList  *Result
}

func NewSites() (s *Sites) {
	s = &Sites{
		conf: conf.NewConfig(),
		VPNList: &Result{
			Vmess:        []string{},
			Vless:        []string{},
			ShadowSocks:  []string{},
			ShadowSocksR: []string{},
			Trojan:       []string{},
		},
	}
	s.siteList = []IVpnSite{
		NewVPNFromGithub(s.conf, s.VPNList),
	}
	return
}

func (that *Sites) Run() {
	for _, s := range that.siteList {
		s.Run()
	}
}
