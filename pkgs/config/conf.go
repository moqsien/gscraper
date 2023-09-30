package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gtea/input"
	"github.com/moqsien/goutils/pkgs/gutils"
)

const (
	GSConfigFileName            string = ".conf_gscraper.json"
	FileInfoName                string = "files_info.json"
	CloudflareRawDomainFileName string = "cloudflare_raw_domains.txt"
	CloudflareDomainFileName    string = "cloudflare_domains.txt"
	NeoboxResultFileName        string = "conf.txt"
	NeoboxResourceLocalDirName  string = "neobox_related"
	GVCResourceLocalDirName     string = "gvc_resources"
	EnableProxyEnvName          string = "ENV_ENABLE_GSER_PROXY"
)

type GVCResourceConfig struct {
	GVCResourceProject string            `json:"gvc_resource_project"`
	GVCResourceDir     string            `json:"gvc_resource_dir"`
	APPUrls            map[string]string `json:"app_url_list"`
	APPUrlOrder        []string          `json:"app_url_order"`
}

type NeoboxResourceConfig struct {
	NeoboxKey             string   `json:"neobox_key"`
	NeoboxResourceProject string   `json:"neobox_resource_project"`
	NeoboxResourceDir     string   `json:"neobox_resource_dir"`
	ProxySubcribeUrlList  []string `json:"subcribe_url_list"`
}

type GSConf struct {
	GVCRConifg    *GVCResourceConfig    `json:"gvc_resource_conifg"`
	NeoboxRConfig *NeoboxResourceConfig `json:"neobox_resource_config"`
	LocalProxy    string                `json:"local_proxy"`
}

func NewGSConf() (gsc *GSConf) {
	gsc = &GSConf{}
	gsc.Initiate()
	return
}

func (that *GSConf) path() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, GSConfigFileName)
}

func (that *GSConf) Initiate() {
	homeDir, _ := os.UserHomeDir()
	configFilePath := filepath.Join(homeDir, GSConfigFileName)
	if ok, _ := gutils.PathIsExist(configFilePath); ok {
		if err := that.Load(); err != nil {
			gprint.PrintError("%+v", err)
			os.RemoveAll(configFilePath)
			os.Exit(1)
		}
	}

	if that.NeoboxRConfig.NeoboxResourceProject == "" || that.GVCRConifg.GVCResourceProject == "" {
		that.SetDefault()
	}

	that.CheckResources()
	if err := that.Save(); err != nil {
		gprint.PrintError("%+v", err)
		os.Exit(1)
	}
}

func (that *GSConf) Load() (err error) {
	var content []byte
	if content, err = os.ReadFile(that.path()); err == nil {
		err = json.Unmarshal(content, that)

	}
	return
}

func (that *GSConf) Save() (err error) {
	var content []byte
	if content, err = json.MarshalIndent(that, "", "    "); err == nil {
		err = os.WriteFile(that.path(), content, os.ModePerm)
	}
	return
}

func (that *GSConf) CheckResources() {
	if that.NeoboxRConfig.NeoboxResourceDir == "" {
		ipt := input.NewInput(input.WithPlaceholder("Neobox Resource Local Dir"), input.WithWidth(100))
		ipt.Run()
		val := ipt.Value()
		if val == "" {
			gprint.PrintError("invalid neobox resource dir")
			os.Exit(1)
		}
		if strings.HasSuffix(val, NeoboxResourceLocalDirName) {
			val = filepath.Dir(val)
		}
		os.MkdirAll(val, os.ModePerm)
		that.NeoboxRConfig.NeoboxResourceDir = filepath.Join(val, NeoboxResourceLocalDirName)
		that.clone(that.NeoboxRConfig.NeoboxResourceDir, that.NeoboxRConfig.NeoboxResourceProject)
	}

	if that.GVCRConifg.GVCResourceDir == "" {
		ipt := input.NewInput(input.WithPlaceholder("GVC Resource Local Dir"), input.WithWidth(100))
		ipt.Run()
		val := ipt.Value()
		if val == "" {
			gprint.PrintError("invalid gvc resource dir")
			os.Exit(1)
		}
		if strings.HasSuffix(val, GVCResourceLocalDirName) {
			val = filepath.Dir(val)
		}
		os.MkdirAll(val, os.ModePerm)
		that.GVCRConifg.GVCResourceDir = filepath.Join(val, GVCResourceLocalDirName)
		that.clone(that.GVCRConifg.GVCResourceDir, that.GVCRConifg.GVCResourceProject)
	}
}

func (that *GSConf) clone(resourceDir, projectUrl string) (err error) {
	if ok, _ := gutils.PathIsExist(resourceDir); ok {
		gprint.PrintInfo("resource local dir: %s already exists", resourceDir)
		return
	}
	git := that.NewGit()
	git.SetWorkDir(filepath.Dir(resourceDir))
	_, err = git.CloneBySSH(projectUrl)
	return
}

func (that *GSConf) SetDefault() {
	that.NeoboxRConfig.NeoboxResourceProject = "git@gitlab.com:moqsien/neobox_related.git"
	that.NeoboxRConfig.ProxySubcribeUrlList = []string{
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
		"https://raw.githubusercontent.com/snakem982/proxypool/main/v2ray3fc8f.txt",
		"https://sub.sharecentre.online/sub",
		"https://getafreenode.com/subscribe/?uuid=D213ED80-199B-4A01-9D62-BBCBA9C16226",
		"http://weoknow.com/data/dayupdate/1/z1.txt",
		"http://weoknow.com/data/dayupdate/2/z1.txt",
		"https://api.subcloud.xyz/sub?target=v2ray&url=https%3A%2F%2Fcdn.jsdelivr.net%2Fgh%2Fzyzmzyz%2Ffree-nodes%40master%2FClash.yml&insert=false",
		"https://api.subcloud.xyz/sub?target=v2ray&url=https%3A%2F%2Fcdn.statically.io%2Fgh%2Fopenrunner%2Fclash-freenode%2Fmain%2Fclash.yaml&insert=false",
		"https://clashnode.com/wp-content/uploads/{year}/{month}/{year}{month}{day}.txt",
		"https://nodefree.org/dy/{year}/{month}/{year}{month}{day}.txt",
		"https://hiclash.com/wp-content/uploads/{year}/{month}/{year}{month}{day}.txt",
		"https://wefound.cc/freenode/{year}/{month}/{year}{month}{day}.txt",
	}

	that.GVCRConifg.GVCResourceProject = "git@gitlab.com:moqsien/gvc_resources.git"
	that.GVCRConifg.APPUrls = map[string]string{
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

	that.GVCRConifg.APPUrlOrder = []string{
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
	that.LocalProxy = "http://localhost:2023"
}
