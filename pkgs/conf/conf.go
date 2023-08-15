package conf

import (
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
		"vlang_linux.zip":        "https://github.com/vlang/v/releases/latest/download/v_linux.zip",
		"vlang_macos.zip":        "https://github.com/vlang/v/releases/latest/download/v_macos.zip",
		"vlang_windows.zip":      "https://github.com/vlang/v/releases/latest/download/v_windows.zip",
		"typst_arm_macos.tar.xz": "https://github.com/typst/typst/releases/latest/download/typst-aarch64-apple-darwin.tar.xz",
		"typst_x64_macos.tar.xz": "https://github.com/typst/typst/releases/latest/download/typst-x86_64-apple-darwin.tar.xz",
		"typst_arm_linux.tar.xz": "https://github.com/typst/typst/releases/latest/download/typst-aarch64-unknown-linux-musl.tar.xz",
		"typst_x64_linux.tar.xz": "https://github.com/typst/typst/releases/latest/download/typst-x86_64-unknown-linux-musl.tar.xz",
		"typst_x64_windows.zip":  "https://github.com/typst/typst/releases/latest/download/typst-x86_64-pc-windows-msvc.zip",
		"nvim_linux64.tar.gz":    "https://github.com/neovim/neovim/releases/download/stable/nvim-linux64.tar.gz",
		"nvim_macos.tar.gz":      "https://github.com/neovim/neovim/releases/download/stable/nvim-macos.tar.gz",
		"nvim_win64.zip":         "https://github.com/neovim/neovim/releases/download/stable/nvim-win64.zip",
		"vcpkg.zip":              "https://github.com/microsoft/vcpkg/archive/refs/heads/master.zip",
		"vcpkg_tool.zip":         "https://github.com/microsoft/vcpkg-tool/archive/refs/heads/main.zip",
		"pyenv_unix.zip":         "https://github.com/pyenv/pyenv/archive/refs/heads/master.zip",
		"pyenv_win.zip":          "https://github.com/pyenv-win/pyenv-win/archive/refs/heads/master.zip",
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
