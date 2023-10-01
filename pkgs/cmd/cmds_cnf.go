package main

import (
	"github.com/moqsien/gscraper/pkgs/config"
	cli "github.com/urfave/cli/v2"
)

func InitCNF() {
	app.Add(&cli.Command{
		Name:    "show-gvc-download-urls",
		Aliases: []string{"show", "sh"},
		Usage:   "Show url list to download for gvc.",
		Action: func(ctx *cli.Context) error {
			cnf := config.NewGSConf()
			cnf.ShowAppUrls()
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "show-neobox-key",
		Aliases: []string{"show-key", "sk"},
		Usage:   "Show neobox key.",
		Action: func(ctx *cli.Context) error {
			cnf := config.NewGSConf()
			cnf.ShowNeoboxKey()
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "generate-neobox-key",
		Aliases: []string{"gen-key", "gk"},
		Usage:   "Generate neobox key.",
		Action: func(ctx *cli.Context) error {
			cnf := config.NewGSConf()
			cnf.SetNeoboxKey()
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "add-gvc-download-urls",
		Aliases: []string{"add-urls", "au"},
		Usage:   "Add url to download for gvc.",
		Action: func(ctx *cli.Context) error {
			sUrl := ctx.Args().First()
			if len(sUrl) > 0 {
				cnf := config.NewGSConf()
				cnf.AddGVCAppUrl(sUrl)
			}
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "delete-gvc-download-urls",
		Aliases: []string{"del-urls", "du"},
		Usage:   "Delete url to download for gvc.",
		Action: func(ctx *cli.Context) error {
			sUrl := ctx.Args().First()
			if len(sUrl) > 0 {
				cnf := config.NewGSConf()
				cnf.DelGVCAppUrl(sUrl)
			}
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "add-neobox-sub-url",
		Aliases: []string{"add-sub", "asub"},
		Usage:   "Add subscribed url for neobox.",
		Action: func(ctx *cli.Context) error {
			sUrl := ctx.Args().First()
			if len(sUrl) > 0 {
				cnf := config.NewGSConf()
				cnf.AddSubscribeUrlForNeobox(sUrl)
			}
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "delete-neobox-sub-url",
		Aliases: []string{"del-sub", "dsub"},
		Usage:   "Delete subscribed url for neobox.",
		Action: func(ctx *cli.Context) error {
			sUrl := ctx.Args().First()
			if len(sUrl) > 0 {
				cnf := config.NewGSConf()
				cnf.DelSubscribeUrlForNeobox(sUrl)
			}
			return nil
		},
	})
}
