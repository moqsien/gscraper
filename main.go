package main

import (
	"fmt"
	"os"
	"sync"

	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/gscraper/pkgs/conf"
	"github.com/moqsien/gscraper/pkgs/download"
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
}

func main() {
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
