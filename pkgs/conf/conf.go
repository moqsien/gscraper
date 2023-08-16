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
	GithubSpeedupUrl   string            `json,koanf:"github_speedup_url"`
	GvcResourceDir     string            `json,koanf:"gvc_resource_dir"`
	GvcResourceProject string            `json,koanf:"gvc_resource_project"`
	UrlList            map[string]string `json,koanf:"url_list"`
	UrlOrder           []string          `json,koanf:"url_order"`
	koanfer            *koanfer.JsonKoanfer
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

func (that *Config) check() {
	if ok, _ := utils.PathIsExist(filepath.Join(that.GvcResourceDir, ".git")); !ok || that.GvcResourceDir == "" {
		tui.PrintError("gvc_resource_dir is required.")
		os.Exit(1)
	}
}

func (that *Config) initiate() {
	if ok, _ := utils.PathIsExist(GScraperConfigPath); !ok {
		that.SetDefault()
		that.koanfer.Save(that)
	}
	if ok, _ := utils.PathIsExist(GScraperConfigPath); ok {
		that.koanfer.Load(that)
	} else {
		tui.PrintError("Cannot find default config files.")
		os.Exit(1)
	}
}

func (that *Config) SetDefault() {
	that.GithubSpeedupUrl = "https://ghproxy.com/"
	that.GvcResourceProject = "git@gitlab.com:moqsien/gvc_resources.git"
	that.UrlList = map[string]string{
		"geoip.db":                        "https://ghproxy.com/?q=https://github.com/lyc8503/sing-box-rules/releases/latest/download/geoip.db",
		"geosite.db":                      "https://ghproxy.com/?q=https://github.com/lyc8503/sing-box-rules/releases/latest/download/geosite.db",
		"geoip.dat":                       "https://ghproxy.com/?q=https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat",
		"geosite.dat":                     "https://ghproxy.com/?q=https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat",
		"protoc_win64.zip":                "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-%s-win64.zip",
		"protoc_linux_x86_64.zip":         "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-%s-linux-x86_64.zip",
		"protoc_linux_aarch_64.zip":       "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-%s-linux-aarch_64.zip",
		"protoc_osx_universal_binary.zip": "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-%s-osx-universal_binary.zip",
		"vlang_linux.zip":                 "https://github.com/vlang/v/releases/latest/download/v_linux.zip",
		"vlang_macos.zip":                 "https://github.com/vlang/v/releases/latest/download/v_macos.zip",
		"vlang_windows.zip":               "https://github.com/vlang/v/releases/latest/download/v_windows.zip",
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
		"geoip.db",
		"geosite.db",
		"geoip.dat",
		"geosite.dat",
		"protoc_win64.zip",
		"protoc_linux_x86_64.zip",
		"protoc_linux_aarch_64.zip",
		"protoc_osx_universal_binary.zip",
		"vlang_linux.zip",
		"vlang_macos.zip",
		"vlang_windows.zip",
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
	that.ReadGvcResourceDir()
}

func (that *Config) ReadGvcResourceDir() {
	item := &tui.InputItem{Title: "GVCResourceProjectDir"}
	input := tui.NewInput([]*tui.InputItem{
		item,
	})
	input.Render()
	that.GvcResourceDir = item.Value
	that.check()

	projectName := "gvc_resources"

	if !strings.Contains(that.GvcResourceDir, projectName) {
		gPath := filepath.Join(that.GvcResourceDir, projectName, ".git")
		if ok, _ := utils.PathIsExist(gPath); ok {
			tui.PrintInfo(fmt.Sprintf("%s already exists.", projectName))
			that.GvcResourceDir = filepath.Join(that.GvcResourceDir, projectName)
			return
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
}

func (that *Config) Add(dUrl string) {
	sList := strings.Split(dUrl, "/")
	filename := sList[len(sList)-1]
	if _, ok := that.UrlList[filename]; !ok {
		that.UrlOrder = append(that.UrlOrder, filename)
	}
	that.UrlList[filename] = dUrl
	that.koanfer.Save(that)
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
	that.koanfer.Save(that)
}

func (that *Config) Show() {
	tui.Cyan(strings.Join(append([]string{"[Files to download]:"}, that.UrlOrder...), "  "))
}
