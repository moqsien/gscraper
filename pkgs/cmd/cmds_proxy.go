package cmd

import (
	"os"

	"github.com/moqsien/gscraper/pkgs/config"
	"github.com/moqsien/gscraper/pkgs/proxy"
	"github.com/moqsien/gscraper/pkgs/proxy/domains"
	"github.com/moqsien/gscraper/pkgs/proxy/proxies"
	cli "github.com/urfave/cli/v2"
)

func InitProxy() {
	app.Add(&cli.Command{
		Name:    "proxy-domains",
		Aliases: []string{"domains", "pm"},
		Usage:   "Get domains for cloudflare edgetunnel.",
		Action: func(ctx *cli.Context) error {
			runner := proxy.NewProxyRunner()
			runner.AddSite(domains.NewCFDomains())
			runner.Run()
			return nil
		},
	})

	var (
		enableGHProxy    bool
		enableLocalProxy bool
	)
	app.Add(&cli.Command{
		Name:    "proxy-rawuri-geoinfo",
		Aliases: []string{"rawuri", "raw"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "use-ghproxy",
				Aliases:     []string{"ghproxy", "g"},
				Destination: &enableGHProxy,
				Usage:       "use ghproxy.com to speedup.",
			},
			&cli.BoolFlag{Name: "use-proxy",
				Aliases:     []string{"proxy", "p"},
				Destination: &enableLocalProxy,
				Usage:       "use local proxy to speedup."},
		},
		Usage: "Get rawURIs and geoinfo for neobox.",
		Action: func(ctx *cli.Context) error {
			runner := proxy.NewProxyRunner()
			runner.AddSite(proxies.NewSubscribers())
			runner.AddSite(proxies.NewWSZiwo())
			runner.AddSite(proxies.NewFreeFQ())
			runner.AddSite(proxies.NewGeoInfo())
			if enableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
			if enableLocalProxy {
				os.Setenv(config.EnableProxyEnvName, "1")
			}
			runner.Run()
			return nil
		},
	})

	app.Add(&cli.Command{
		Name:    "proxy-rawuri-only",
		Aliases: []string{"ruo", "ro"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "use-ghproxy",
				Aliases:     []string{"ghproxy", "g"},
				Destination: &enableGHProxy,
				Usage:       "use ghproxy.com to speedup.",
			},
			&cli.BoolFlag{Name: "use-proxy",
				Aliases:     []string{"proxy", "p"},
				Destination: &enableLocalProxy,
				Usage:       "use local proxy to speedup."},
		},
		Usage: "Get rawURIs and geoinfo for neobox.",
		Action: func(ctx *cli.Context) error {
			runner := proxy.NewProxyRunner()
			runner.AddSite(proxies.NewSubscribers())
			runner.AddSite(proxies.NewWSZiwo())
			runner.AddSite(proxies.NewFreeFQ())
			if enableGHProxy {
				os.Setenv(config.EnableGithubSpeedupEnvName, "1")
			}
			if enableLocalProxy {
				os.Setenv(config.EnableProxyEnvName, "1")
			}
			runner.Run()
			return nil
		},
	})
}
