package main

import (
	"github.com/moqsien/gscraper/pkgs/gapps"
	cli "github.com/urfave/cli/v2"
)

func IniteGVC() {
	app.Add(&cli.Command{
		Name:    "download-apps-for-gvc",
		Aliases: []string{"dapps", "dag"},
		Usage:   "Download gvc files.",
		Action: func(ctx *cli.Context) error {
			fNames := ctx.Args().Slice()
			d := gapps.NewDownloader()
			d.Start(fNames...)
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "download-gvc",
		Aliases: []string{"downgvc", "dg"},
		Usage:   "Download gvc files.",
		Action: func(ctx *cli.Context) error {
			d := gapps.NewDownloader()
			fNames := []string{
				"gvc_darwin-amd64.zip",
				"gvc_darwin-arm64.zip",
				"gvc_linux-amd64.zip",
				"gvc_linux-arm64.zip",
				"gvc_windows-amd64.zip",
				"gvc_windows-arm64.zip",
			}
			d.Start(fNames...)
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "download-other",
		Aliases: []string{"downother", "do"},
		Usage:   "Download files except gvc.",
		Action: func(ctx *cli.Context) error {
			d := gapps.NewDownloader()
			fNames := []string{
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
			d.Start(fNames...)
			return nil
		},
	})
}
