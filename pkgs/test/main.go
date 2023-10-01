package main

import (
	"os"

	"github.com/moqsien/gscraper/pkgs/config"
	"github.com/moqsien/gscraper/pkgs/proxy"
	"github.com/moqsien/gscraper/pkgs/proxy/domains"
	"github.com/moqsien/gscraper/pkgs/proxy/proxies"
)

func main() {
	// os.Setenv(config.EnableGithubSpeedupEnvName, "1")
	// sub := proxies.NewSubscribers()
	// sub.SetHandler(func(result []string) {
	// 	gprint.PrintInfo("find rawUris: %d", len(result))
	// })
	// sub.Run()

	// fproxy := proxies.NewWSZiwo()
	// fproxy.SetHandler(func(result []string) {
	// 	gprint.PrintInfo("find rawUris: %d", len(result))
	// })
	// fproxy.Run()

	// os.Setenv(config.EnableProxyEnvName, "1")
	// fq := proxies.NewFreeFQ()
	// fq.SetHandler(func(s []string) {
	// 	fmt.Println(s)
	// 	fmt.Println(len(s))
	// })
	// fq.Run()

	runner := proxy.NewProxyRunner()
	runner.AddSite(proxies.NewSubscribers())
	runner.AddSite(proxies.NewWSZiwo())
	runner.AddSite(proxies.NewFreeFQ())
	runner.AddSite(proxies.NewGeoInfo())
	runner.AddSite(domains.NewCFDomains())
	os.Setenv(config.EnableGithubSpeedupEnvName, "1")
	os.Setenv(config.EnableProxyEnvName, "1")

	runner.Run()
}
