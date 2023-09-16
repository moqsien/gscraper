package main

import (
	"os"
	"strings"

	"github.com/moqsien/gscraper/pkgs/cmd"
	"github.com/moqsien/gscraper/pkgs/conf"
	"github.com/moqsien/gscraper/pkgs/domain"
)

func main() {
	p, _ := os.Executable()

	if strings.HasSuffix(p, "gscraper") || strings.HasSuffix(p, "gscraper.exe") {
		cmd.StartApp()
	} else {
		// cnf := conf.NewConfig()
		// cnf.SetDefault()
		// r := &sites.Result{
		// 	Vmess:        []string{},
		// 	Vless:        []string{},
		// 	ShadowSocks:  []string{},
		// 	ShadowSocksR: []string{},
		// 	Trojan:       []string{},
		// }
		// vg := sites.NewVPNFromGithub(cnf, r)
		// vg.Run()
		// fmt.Println(vg.VPNList)
		// u := utils.NewUUID()
		// fmt.Println(strings.ToUpper(u.String()))
		cnf := conf.NewConfig()
		cdomain := domain.NewCFlareDomain(cnf)
		cdomain.Run()
	}
}
