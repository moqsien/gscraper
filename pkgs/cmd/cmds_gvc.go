package main

import (
	"os"

	"github.com/moqsien/gscraper/pkgs/config"
	"github.com/moqsien/gscraper/pkgs/gapps"
	cli "github.com/urfave/cli/v2"
)

func IniteGVC() {
	var disableGHProxy bool
	app.Add(&cli.Command{
		Name:    "download-apps-for-gvc",
		Aliases: []string{"dapps", "dag"},
		Usage:   "Download gvc files.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable-ghproxy",
				Aliases:     []string{"dg", "d"},
				Destination: &disableGHProxy,
				Usage:       "disable using ghproxy.com to speedup.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !disableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}

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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable-ghproxy",
				Aliases:     []string{"dg", "d"},
				Destination: &disableGHProxy,
				Usage:       "disable using ghproxy.com to speedup.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !disableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable-ghproxy",
				Aliases:     []string{"dg", "d"},
				Destination: &disableGHProxy,
				Usage:       "disable using ghproxy.com to speedup.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !disableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
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

	app.Add(&cli.Command{
		Name:    "download-master",
		Aliases: []string{"downmaster", "dm"},
		Usage:   "Download github master files.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable-ghproxy",
				Aliases:     []string{"dg", "d"},
				Destination: &disableGHProxy,
				Usage:       "disable using ghproxy.com to speedup.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !disableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
			d := gapps.NewDownloader()
			fNames := []string{
				"vcpkg.zip",
				"vcpkg_tool.zip",
				"pyenv_unix.zip",
				"pyenv_win.zip",
			}
			d.Start(fNames...)
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "download-nvim",
		Aliases: []string{"downvim", "dn"},
		Usage:   "Download nvim files.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable-ghproxy",
				Aliases:     []string{"dg", "d"},
				Destination: &disableGHProxy,
				Usage:       "disable using ghproxy.com to speedup.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !disableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
			d := gapps.NewDownloader()
			fNames := []string{
				"nvim_linux64.tar.gz",
				"nvim_macos.tar.gz",
				"nvim_win64.zip",
			}
			d.Start(fNames...)
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "download-typst",
		Aliases: []string{"downtypst", "dt"},
		Usage:   "Download typst files.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable-ghproxy",
				Aliases:     []string{"dg", "d"},
				Destination: &disableGHProxy,
				Usage:       "disable using ghproxy.com to speedup.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !disableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
			d := gapps.NewDownloader()
			fNames := []string{
				"typst_arm_macos.tar.xz",
				"typst_x64_macos.tar.xz",
				"typst_arm_linux.tar.xz",
				"typst_x64_linux.tar.xz",
				"typst_x64_windows.zip",
			}
			d.Start(fNames...)
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "download-vlang",
		Aliases: []string{"downvlang", "dv"},
		Usage:   "Download vlang files.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable-ghproxy",
				Aliases:     []string{"dg", "d"},
				Destination: &disableGHProxy,
				Usage:       "disable using ghproxy.com to speedup.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !disableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
			d := gapps.NewDownloader()
			fNames := []string{
				"vlang_linux.zip",
				"vlang_macos.zip",
				"vlang_windows.zip",
				"v_analyzer_darwin_arm64.zip",
				"v_analyzer_darwin_x86_64.zip",
				"v_analyzer_linux_x86_64.zip",
				"v_analyzer_windows_x86_64.zip",
			}
			d.Start(fNames...)
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "download-protobuf",
		Aliases: []string{"downprotobuf", "dpb"},
		Usage:   "Download protobuf files.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "disable-ghproxy",
				Aliases:     []string{"dg", "d"},
				Destination: &disableGHProxy,
				Usage:       "disable using ghproxy.com to speedup.",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !disableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
			d := gapps.NewDownloader()
			fNames := []string{
				"protoc_win64.zip",
				"protoc_linux_x86_64.zip",
				"protoc_linux_aarch_64.zip",
				"protoc_osx_universal_binary.zip",
			}
			d.Start(fNames...)
			return nil
		},
	})
}
