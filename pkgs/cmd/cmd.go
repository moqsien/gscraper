package cmd

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/gscraper/pkgs/conf"
	"github.com/moqsien/gscraper/pkgs/download"
	"github.com/moqsien/gscraper/pkgs/sites"
	cli "github.com/urfave/cli/v2"
)

type App struct {
	cmd *cli.App
}

func New() *App {
	return &App{
		cmd: &cli.App{
			Usage:       "gscraper <Command> <SubCommand>...",
			Description: "gscraper, download files from github for gvc.",
			Commands:    []*cli.Command{},
		},
	}
}

func (that *App) Add(command *cli.Command) {
	that.cmd.Commands = append(that.cmd.Commands, command)
}

func (that *App) Run() {
	that.cmd.Run(os.Args)
}

var app *App

func init() {
	app = New()
	app.Add(&cli.Command{
		Name:    "show",
		Aliases: []string{"sh"},
		Usage:   "Show url list to download.",
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			cnf.Show()
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add new download urls[ gscraper add xxx yyy zzz ss].",
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			cnf := conf.NewConfig()
			for _, arg := range args {
				if utils.VerifyUrls(arg) {
					cnf.Add(arg)
				}
			}
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "remove",
		Aliases: []string{"rm", "r"},
		Usage:   "Remove download url.[ gscraper rm xxx ]",
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			arg := ctx.Args().First()
			cnf.Remove(arg)
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "reset",
		Aliases: []string{"rs"},
		Usage:   "Reset config file.",
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			cnf.SetDefault()
			cnf.Save()
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "download",
		Aliases: []string{"down", "d"},
		Usage:   "Download files.",
		Action: func(ctx *cli.Context) error {
			d := download.NewDownloader()
			args := ctx.Args().Slice()
			d.Start(args...)
			return nil
		},
	})

	vpnCli := &cli.Command{
		Name:        "vpn",
		Aliases:     []string{"v"},
		Usage:       "free vpn related.",
		Subcommands: []*cli.Command{},
	}

	vshow := &cli.Command{
		Name:    "show",
		Usage:   "Show subscribers from github.",
		Aliases: []string{"sh", "s"},
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			cnf.ShowGithubVPNSubscriber()
			return nil
		},
	}
	vpnCli.Subcommands = append(vpnCli.Subcommands, vshow)

	vadd := &cli.Command{
		Name:    "add",
		Usage:   "Add subscribers from github.",
		Aliases: []string{"ad", "a"},
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			sUrl := ctx.Args().First()
			cnf.AddGithubVpnSubscriber(sUrl)
			return nil
		},
	}
	vpnCli.Subcommands = append(vpnCli.Subcommands, vadd)

	vrm := &cli.Command{
		Name:    "remove",
		Usage:   "Remove subscribers from github.",
		Aliases: []string{"rm", "r"},
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			if idx, err := strconv.Atoi(ctx.Args().First()); err == nil {
				cnf.RemoveGithubVPNSubscriber(idx)
			}
			return nil
		},
	}
	vpnCli.Subcommands = append(vpnCli.Subcommands, vrm)

	vproxy := &cli.Command{
		Name:    "proxy",
		Usage:   "Set local proxy for gscraper.",
		Aliases: []string{"pxy", "p"},
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			pxy := ctx.Args().First()
			cnf.SetLocalProxy(pxy)
			return nil
		},
	}
	vpnCli.Subcommands = append(vpnCli.Subcommands, vproxy)

	vnKey := &cli.Command{
		Name:    "neobox-key",
		Usage:   "Reset neobox-key for gscraper.",
		Aliases: []string{"nk", "key"},
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			cnf.ResetNeoboxKey()
			return nil
		},
	}
	vpnCli.Subcommands = append(vpnCli.Subcommands, vnKey)

	vrun := &cli.Command{
		Name:    "get-free-vpn",
		Usage:   "Get free vpns.",
		Aliases: []string{"gfv", "free"},
		Action: func(ctx *cli.Context) error {
			s := sites.NewSites()
			s.Run()
			return nil
		},
	}
	vpnCli.Subcommands = append(vpnCli.Subcommands, vrun)

	vshownk := &cli.Command{
		Name:    "show-key",
		Usage:   "Show neobox-key.",
		Aliases: []string{"sk", "skey"},
		Action: func(ctx *cli.Context) error {
			cnf := conf.NewConfig()
			cnf.ShowNeoboxKey()
			return nil
		},
	}
	vpnCli.Subcommands = append(vpnCli.Subcommands, vshownk)
	app.Add(vpnCli)
}

func StartApp() {
	sig := &utils.CtrlCSignal{}
	w := &sync.WaitGroup{}
	sig.RegisterSweeper(func() error {
		w.Add(1)
		d := download.NewDownloader()
		download.GLOBAL_TO_EXIST = true
		_, ok := <-download.WaitToSweepSig
		if !ok {
			fmt.Println("remove temporary files...")
			d.RemoveTempDir()
			fmt.Println("push downloaded file to remote repository...")
			d.GitPush()
		}
		w.Done()
		return nil
	})
	sig.ListenSignal()
	app.Run()
	w.Wait()
}
