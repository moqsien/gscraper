package sites

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/moqsien/goutils/pkgs/crypt"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/gscraper/pkgs/conf"
	"github.com/moqsien/vpnparser/pkgs/outbound"
)

/*
Gather free proxies for neobox/gvc-neobox.
*/

const (
	JSON_FILE_NAME = "free_vpn_list.json"
)

type IVpnSite interface {
	Parse([]byte)
	Run()
}

type Sites struct {
	conf     *conf.Config
	siteList []IVpnSite
	VPNList  *outbound.Result `json:"vpn_list"`
	path     string
}

func NewSites() (s *Sites) {
	s = &Sites{
		conf:     conf.NewConfig(),
		VPNList:  outbound.NewResult(),
		siteList: []IVpnSite{},
	}
	s.RegisterSite(NewVPNFromGithub(s.conf, s.VPNList))
	s.path = filepath.Join(s.conf.GvcResourceDir, JSON_FILE_NAME)
	return
}

func (that *Sites) RegisterSite(ivs IVpnSite) {
	if ivs == nil {
		return
	}
	that.siteList = append(that.siteList, ivs)
}

func (that *Sites) setGitIgnore() {
	gPath := filepath.Join(that.conf.GvcResourceDir, ".gitignore")
	if ok, _ := utils.PathIsExist(gPath); !ok {
		os.WriteFile(gPath, []byte(JSON_FILE_NAME), os.ModePerm)
	} else {
		content, _ := os.ReadFile(gPath)
		if !strings.Contains(string(content), JSON_FILE_NAME) {
			content = append(content, []byte(fmt.Sprintf("\n%s", JSON_FILE_NAME))...)
			os.WriteFile(gPath, content, os.ModePerm)
		}
	}
}

func (that *Sites) Save() {
	tui.PrintInfo(fmt.Sprintf("vmess[%d], vless[%d], ss[%d], ssr[%d], trojan[%d]",
		len(that.VPNList.Vmess),
		len(that.VPNList.Vless),
		len(that.VPNList.ShadowSocks),
		len(that.VPNList.ShadowSocksR),
		len(that.VPNList.Trojan)))

	that.VPNList.UpdateAt = time.Now().Format("2006-01-02 15:04:05")
	that.VPNList.VmessTotal = len(that.VPNList.Vmess)
	that.VPNList.VlessTotal = len(that.VPNList.Vless)
	that.VPNList.SSTotal = len(that.VPNList.ShadowSocks)
	that.VPNList.SSRTotal = len(that.VPNList.ShadowSocksR)
	that.VPNList.TrojanTotal = len(that.VPNList.Trojan)

	that.VPNList.Save(that.path)
	that.setGitIgnore()
}

func (that *Sites) push() {
	cmdExe := "git"
	if runtime.GOOS == "windows" {
		cmdExe = "git.exe"
	}
	cmd := exec.Command(cmdExe, "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Dir = that.conf.GvcResourceDir
	err := cmd.Run()
	if err != nil {
		tui.PrintError(err)
		return
	}

	cmd = exec.Command(cmdExe, "commit", "-m", "update")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Dir = that.conf.GvcResourceDir
	err = cmd.Run()
	if err != nil {
		tui.PrintError(err)
		return
	}

	cmd = exec.Command(cmdExe, "push")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Dir = that.conf.GvcResourceDir
	err = cmd.Run()
	if err != nil {
		tui.PrintError(err)
		return
	}
}

func (that *Sites) GitPush() {
	that.Save()
	fPath := filepath.Join(that.conf.GvcResourceDir, that.conf.NeoboxResultFile)
	content, _ := os.ReadFile(filepath.Join(that.conf.GvcResourceDir, JSON_FILE_NAME))
	if len(content) > 0 {
		cryp := crypt.NewCrptWithKey([]byte(that.conf.NeoboxKey))
		if data, err := cryp.AesEncrypt(content); err == nil {
			os.WriteFile(fPath, data, 0777)
			that.push()
		}
	}
}

func (that *Sites) Run() {
	for _, s := range that.siteList {
		s.Run()
	}
	that.GitPush()
}
