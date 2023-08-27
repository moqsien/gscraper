package main

import (
	"fmt"
	"os"
	"strings"

	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/gscraper/pkgs/cmd"
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
		u := utils.NewUUID()
		fmt.Println(strings.ToUpper(u.String()))
	}
}
