package conf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/goutils/pkgs/koanfer"
)

var (
	GScraperConfigPath string = func() string {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "gscraper_conf.json")
	}()
)

type Config struct {
	GithubSpeedupUrl     string            `json,koanf:"github_speedup_url"`
	GvcResourceDir       string            `json,koanf:"gvc_resource_dir"`
	GvcResourceProject   string            `json,koanf:"gvc_resource_project"`
	UrlList              map[string]string `json,koanf:"url_list"`
	UrlOrder             []string          `json,koanf:"url_order"`
	GithubVPNSubscriber  []string          `json,koanf:"github_vpn_subscriber"`
	CloudflareDomainFile string            `json,koanf:"cloudflare_domain_filename"`
	CloudflareDomains    []string          `json,koanf:"cloudflare_domains"`
	LocalProxy           string            `json,koanf:"local_proxy"`
	NeoboxKey            string            `json,koanf:"neobox_key"`
	NeoboxResultFile     string            `json,koanf:"neobox_result_file"`
	koanfer              *koanfer.JsonKoanfer
}

func NewConfig() *Config {
	kfer, _ := koanfer.NewKoanfer(GScraperConfigPath)
	cfg := &Config{
		UrlList: map[string]string{},
		koanfer: kfer,
	}
	cfg.initiate()
	cfg.check()
	return cfg
}

func (that *Config) Save() {
	that.koanfer.Save(that)
}

func (that *Config) check() {
	if ok, _ := utils.PathIsExist(filepath.Join(that.GvcResourceDir, ".git")); that.GvcResourceDir == "" || !ok {
		tui.PrintWarning("gvc_resource_dir is required. [gvc_resuorces_dir]: ", that.GvcResourceDir)
		that.ReadGvcResourceDir()
		os.Exit(1)
	}
}

func (that *Config) initiate() {
	if ok, _ := utils.PathIsExist(GScraperConfigPath); !ok {
		that.SetDefault()
		that.Save()
	}
	if ok, _ := utils.PathIsExist(GScraperConfigPath); ok {
		that.koanfer.Load(that)
	} else {
		tui.PrintError("Cannot find default config files.")
		os.Exit(1)
	}
}

func (that *Config) SetDefault() {
	that.LocalProxy = "http://localhost:2023"
	that.NeoboxKey = "IYj8oCV1Nly9aTTN"
	that.GithubSpeedupUrl = "https://ghproxy.com/"
	that.GvcResourceProject = "git@gitlab.com:moqsien/gvc_resources.git"
	that.UrlList = map[string]string{
		"gsudo_portable.zip":              "https://github.com/gerardog/gsudo/releases//latest/download/gsudo.portable.zip",
		"gvc_darwin-amd64.zip":            "https://github.com/moqsien/gvc/releases/latest/download/gvc_darwin-amd64.zip",
		"gvc_darwin-arm64.zip":            "https://github.com/moqsien/gvc/releases/latest/download/gvc_darwin-arm64.zip",
		"gvc_linux-amd64.zip":             "https://github.com/moqsien/gvc/releases/latest/download/gvc_linux-amd64.zip",
		"gvc_linux-arm64.zip":             "https://github.com/moqsien/gvc/releases/latest/download/gvc_linux-arm64.zip",
		"gvc_windows-amd64.zip":           "https://github.com/moqsien/gvc/releases/latest/download/gvc_windows-amd64.zip",
		"gvc_windows-arm64.zip":           "https://github.com/moqsien/gvc/releases/latest/download/gvc_windows-arm64.zip",
		"geoip.db":                        "https://github.com/lyc8503/sing-box-rules/releases/latest/download/geoip.db",
		"geosite.db":                      "https://github.com/lyc8503/sing-box-rules/releases/latest/download/geosite.db",
		"geoip.dat":                       "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat",
		"geosite.dat":                     "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat",
		"protoc_win64.zip":                "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-%s-win64.zip",
		"protoc_linux_x86_64.zip":         "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-%s-linux-x86_64.zip",
		"protoc_linux_aarch_64.zip":       "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-%s-linux-aarch_64.zip",
		"protoc_osx_universal_binary.zip": "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-%s-osx-universal_binary.zip",
		"vlang_linux.zip":                 "https://github.com/vlang/v/releases/latest/download/v_linux.zip",
		"vlang_macos.zip":                 "https://github.com/vlang/v/releases/latest/download/v_macos.zip",
		"vlang_windows.zip":               "https://github.com/vlang/v/releases/latest/download/v_windows.zip",
		"v_analyzer_darwin_arm64.zip":     "https://github.com/v-analyzer/v-analyzer/releases/latest/download/v-analyzer-darwin-arm64.zip",
		"v_analyzer_darwin_x86_64.zip":    "https://github.com/v-analyzer/v-analyzer/releases/latest/download/v-analyzer-darwin-x86_64.zip",
		"v_analyzer_linux_x86_64.zip":     "https://github.com/v-analyzer/v-analyzer/releases/latest/download/v-analyzer-linux-x86_64.zip",
		"v_analyzer_windows_x86_64.zip":   "https://github.com/v-analyzer/v-analyzer/releases/latest/download/v-analyzer-windows-x86_64.zip",
		"typst_arm_macos.tar.xz":          "https://github.com/typst/typst/releases/latest/download/typst-aarch64-apple-darwin.tar.xz",
		"typst_x64_macos.tar.xz":          "https://github.com/typst/typst/releases/latest/download/typst-x86_64-apple-darwin.tar.xz",
		"typst_arm_linux.tar.xz":          "https://github.com/typst/typst/releases/latest/download/typst-aarch64-unknown-linux-musl.tar.xz",
		"typst_x64_linux.tar.xz":          "https://github.com/typst/typst/releases/latest/download/typst-x86_64-unknown-linux-musl.tar.xz",
		"typst_x64_windows.zip":           "https://github.com/typst/typst/releases/latest/download/typst-x86_64-pc-windows-msvc.zip",
		"nvim_linux64.tar.gz":             "https://github.com/neovim/neovim/releases/download/stable/nvim-linux64.tar.gz",
		"nvim_macos.tar.gz":               "https://github.com/neovim/neovim/releases/download/stable/nvim-macos.tar.gz",
		"nvim_win64.zip":                  "https://github.com/neovim/neovim/releases/download/stable/nvim-win64.zip",
		"vcpkg.zip":                       "https://github.com/microsoft/vcpkg/archive/refs/heads/master.zip",
		"vcpkg_tool.zip":                  "https://github.com/microsoft/vcpkg-tool/archive/refs/heads/main.zip",
		"pyenv_unix.zip":                  "https://github.com/pyenv/pyenv/archive/refs/heads/master.zip",
		"pyenv_win.zip":                   "https://github.com/pyenv-win/pyenv-win/archive/refs/heads/master.zip",
	}
	that.UrlOrder = []string{
		"gvc_darwin-amd64.zip",
		"gvc_darwin-arm64.zip",
		"gvc_linux-amd64.zip",
		"gvc_linux-arm64.zip",
		"gvc_windows-amd64.zip",
		"gvc_windows-arm64.zip",
		"geoip.db",
		"geosite.db",
		"geoip.dat",
		"geosite.dat",
		"gsudo_portable.zip",
		"protoc_win64.zip",
		"protoc_linux_x86_64.zip",
		"protoc_linux_aarch_64.zip",
		"protoc_osx_universal_binary.zip",
		"vlang_linux.zip",
		"vlang_macos.zip",
		"vlang_windows.zip",
		"v_analyzer_darwin_arm64.zip",
		"v_analyzer_darwin_x86_64.zip",
		"v_analyzer_linux_x86_64.zip",
		"v_analyzer_windows_x86_64.zip",
		"typst_arm_macos.tar.xz",
		"typst_x64_macos.tar.xz",
		"typst_arm_linux.tar.xz",
		"typst_x64_linux.tar.xz",
		"typst_x64_windows.zip",
		"nvim_linux64.tar.gz",
		"nvim_macos.tar.gz",
		"nvim_win64.zip",
		"vcpkg.zip",
		"vcpkg_tool.zip",
		"pyenv_unix.zip",
		"pyenv_win.zip",
	}
	that.GithubVPNSubscriber = []string{
		"https://raw.githubusercontent.com/peasoft/NoMoreWalls/master/list.txt",
		"https://raw.githubusercontent.com/ZywChannel/free/main/sub",
		"https://raw.githubusercontent.com/ermaozi01/free_clash_vpn/main/subscribe/v2ray.txt",
		"https://raw.githubusercontent.com/Pawdroid/Free-servers/main/sub",
		"https://raw.githubusercontent.com/freefq/free/master/v2",
		"https://raw.githubusercontent.com/mfuu/v2ray/master/v2ray",
		"https://raw.githubusercontent.com/ssrsub/ssr/master/ss-sub",
		"https://raw.githubusercontent.com/ssrsub/ssr/master/V2Ray",
		"https://raw.githubusercontent.com/ermaozi/get_subscribe/main/subscribe/v2ray.txt",
		"https://raw.githubusercontent.com/tbbatbb/Proxy/master/dist/v2ray.config.txt",
		"https://raw.githubusercontent.com/vveg26/get_proxy/main/dist/v2ray.config.txt",
		"https://raw.githubusercontent.com/baip01/yhkj/main/v2ray",
		"https://raw.githubusercontent.com/aiboboxx/v2rayfree/main/v2",
		"https://raw.githubusercontent.com/ts-sf/fly/main/v2",
		"https://raw.githubusercontent.com/free18/v2ray/main/v2ray.txt",
		"https://raw.githubusercontent.com/Leon406/SubCrawler/main/sub/share/vless",
		"https://raw.githubusercontent.com/Leon406/SubCrawler/main/sub/share/ss",
		"https://raw.githubusercontent.com/Leon406/SubCrawler/main/sub/share/ssr",
		"https://raw.githubusercontent.com/Leon406/SubCrawler/main/sub/share/all3",
		"https://raw.githubusercontent.com/Leon406/SubCrawler/main/sub/share/v2",
		"https://sub.sharecentre.online/sub",
		"https://getafreenode.com/subscribe/?uuid=D213ED80-199B-4A01-9D62-BBCBA9C16226",
		"https://wanshanziwo.eu.org/vmess/sub?rand=lNdLVPVC&c=US",
		"https://wanshanziwo.eu.org/vless/sub?rand=lNdLVPVC&c=US",
		"https://wanshanziwo.eu.org/trojan/sub?rand=lNdLVPVC",
		"https://wanshanziwo.eu.org/ss/sub?rand=lNdLVPVC&c=US",
		"https://wanshanziwo.eu.org/ssr/sub?rand=lNdLVPVC",
		"http://weoknow.com/data/dayupdate/1/z1.txt",
		"http://weoknow.com/data/dayupdate/2/z1.txt",
		"https://api.subcloud.xyz/sub?target=v2ray&url=https%3A%2F%2Fcdn.jsdelivr.net%2Fgh%2Fzyzmzyz%2Ffree-nodes%40master%2FClash.yml&insert=false",
		"https://api.subcloud.xyz/sub?target=v2ray&url=https%3A%2F%2Fcdn.statically.io%2Fgh%2Fopenrunner%2Fclash-freenode%2Fmain%2Fclash.yaml&insert=false",
		"https://clashnode.com/wp-content/uploads/{year}/{month}/{year}{month}{day}.txt",
		"https://nodefree.org/dy/{year}/{month}/{year}{month}{day}.txt",
		"https://hiclash.com/wp-content/uploads/{year}/{month}/{year}{month}{day}.txt",
		"https://wefound.cc/freenode/{year}/{month}/{year}{month}{day}.txt",
	}
	that.NeoboxResultFile = "conf.txt"

	that.CloudflareDomainFile = "cloudflare_domains.txt"
	that.CloudflareDomains = []string{
		"time.cloudflare.com",
		// "shopify.com",
		// "time.is",
		"icook.hk",
		"icook.tw",
		"ip.sb",
		"japan.com",
		"malaysia.com",
		"russia.com",
		"singapore.com",
		"www.visa.com",
		"www.visa.com.sg",
		"www.visa.com.hk",
		"www.visa.com.tw",
		"www.visa.co.jp",
		"www.visakorea.com",
		"www.gco.gov.qa",
		// "www.gov.se",
		"www.gov.ua",
		"www.digitalocean.com",
		"www.csgo.com",
		"www.shopify.com",
		"www.whoer.net",
		"www.whatismyip.com",
		"www.ipget.net",
		"www.hugedornains.com",
		"www.udacity.com",
		"www.4chan.org",
		"www.okcupid.com",
		"www.glassdoor.com",
		"www.udemy.com",
		"www.baipiao.eu.org",
		// "cdn.anycast.eu.org",
		// "cdn-all.xn--b6gac.eu.org",
		// "cdn-b100.xn--b6gac.eu.org",
		// "cdn.xn--b6gac.eu.org",
		// "edgetunnel.anycast.eu.org",
	}
	that.Save()
	that.ReadGvcResourceDir()
}

func (that *Config) ReadGvcResourceDir() {
	item := &tui.InputItem{Title: "GVCResourceProjectDir"}
	input := tui.NewInput([]*tui.InputItem{
		item,
	})
	input.Render()
	if item.Value == "" {
		tui.PrintInfo("[gvc_resource dir]: ", that.GvcResourceDir)
	} else {
		that.GvcResourceDir = item.Value
	}
	projectName := "gvc_resources"

	if !strings.Contains(that.GvcResourceDir, projectName) {
		gPath := filepath.Join(that.GvcResourceDir, projectName, ".git")
		if ok, _ := utils.PathIsExist(gPath); ok {
			tui.PrintInfo(fmt.Sprintf("%s already exists.", projectName))
			that.GvcResourceDir = filepath.Join(that.GvcResourceDir, projectName)
			return
		} else {
			os.MkdirAll(that.GvcResourceDir, 0777)
		}
		cmdName := "git"
		if runtime.GOOS == "windows" {
			cmdName = "git.exe"
		}
		cmd := exec.Command(cmdName, "clone", that.GvcResourceProject)
		cmd.Dir = that.GvcResourceDir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			tui.PrintError(err)
			os.Exit(1)
		}
		that.GvcResourceDir = filepath.Join(that.GvcResourceDir, projectName)
	}
	that.Save()
}

func (that *Config) getName(dUrl string) string {
	sList := strings.Split(dUrl, "/")
	filename := sList[len(sList)-1]
	if strings.Contains(filename, "master") || strings.Contains(filename, "main") {
		sList = strings.Split(dUrl, "github.com/")
		if len(sList) == 2 {
			s := sList[1]
			sList = strings.Split(s, "/")
			if len(sList) > 1 {
				return fmt.Sprintf("%s.zip", sList[1])
			}
		}
	}
	return filename
}

func (that *Config) Add(dUrl string) {
	filename := that.getName(dUrl)
	that.UrlList[filename] = dUrl
	that.Save()
}

func (that *Config) Remove(filename string) {
	for idx, name := range that.UrlOrder {
		if name == filename {
			if idx != len(that.UrlOrder)-1 {
				that.UrlOrder = append(that.UrlOrder[:idx], that.UrlOrder[idx+1:]...)
			} else {
				that.UrlOrder = that.UrlOrder[:idx]
			}
			break
		}
	}
	delete(that.UrlList, filename)
	os.RemoveAll(filepath.Join(that.GvcResourceDir, filename))
	that.Save()
}

func (that *Config) Show() {
	tui.Cyan(strings.Join(append([]string{"[Files to download]:"}, that.UrlOrder...), "  "))
}

func (that *Config) AddGithubVpnSubscriber(sUrl string) {
	flag := false
	for _, v := range that.GithubVPNSubscriber {
		if v == sUrl {
			flag = true
			break
		}
	}
	if !flag && utils.VerifyUrls(sUrl) {
		that.GithubVPNSubscriber = append(that.GithubVPNSubscriber, sUrl)
		that.Save()
	}
}

func (that *Config) RemoveGithubVPNSubscriber(index int) {
	if index < 0 || index >= len(that.GithubVPNSubscriber) {
		return
	}
	if index == len(that.GithubVPNSubscriber)-1 {
		that.GithubVPNSubscriber = that.GithubVPNSubscriber[:index]
	} else {
		that.GithubVPNSubscriber = append(that.GithubVPNSubscriber[:index], that.GithubVPNSubscriber[index+1:]...)
	}
	that.Save()
}

func (that *Config) ShowGithubVPNSubscriber() {
	for idx, sUrl := range that.GithubVPNSubscriber {
		fmt.Printf("%v. %s\n", idx, sUrl)
	}
}

func (that *Config) SetLocalProxy(pxy string) {
	that.LocalProxy = pxy
	that.Save()
}

func (that *Config) ResetNeoboxKey() {
	that.NeoboxKey = utils.RandomString(16)
	that.Save()
}

func (that *Config) ShowNeoboxKey() {
	tui.PrintInfo(fmt.Sprintf("neobox-key: %s\n", that.NeoboxKey))
}

func (that *Config) AddCFlareDomain(cDomain string) {
	that.koanfer.Load(that)
	that.CloudflareDomains = append(that.CloudflareDomains, cDomain)
	that.Save()
}

func (that *Config) RemoveCFlareDomain(cDomain string) {
	that.koanfer.Load(that)
	for idx, domain := range that.CloudflareDomains {
		if domain == cDomain {
			that.CloudflareDomains = append(that.CloudflareDomains[:idx], that.CloudflareDomains[idx+1:]...)
			break
		}
	}
	that.Save()
}
