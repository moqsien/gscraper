package main

import (
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gscraper/pkgs/proxy/proxies"
)

func main() {
	// os.Setenv(config.EnableGithubSpeedupEnvName, "1")
	// sub := proxies.NewSubscribers()
	// sub.SetHandler(func(result []string) {
	// 	gprint.PrintInfo("find rawUris: %d", len(result))
	// })
	// sub.Run()

	fproxy := proxies.NewWSZiwo()
	fproxy.SetHandler(func(result []string) {
		gprint.PrintInfo("find rawUris: %d", len(result))
	})
	fproxy.Run()
}
