package main

import (
	"fmt"
	"os"
	"sync"

	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/gscraper/pkgs/gapps"
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

func StartApp() {
	sig := &utils.CtrlCSignal{}
	w := &sync.WaitGroup{}
	sig.RegisterSweeper(func() error {
		w.Add(1)
		d := gapps.NewDownloader()
		gapps.GLOBAL_TO_EXIST = true
		_, ok := <-gapps.WaitToSweepSig
		if !ok {
			fmt.Println("remove temporary files...")
			d.RemoveTempDir()
			fmt.Println("push downloaded file to remote repository...")
			d.Push()
		}
		w.Done()
		return nil
	})
	sig.ListenSignal()
	app.Run()
	w.Wait()
}

func main() {
	StartApp()
}
