package config

import (
	"os"

	"github.com/moqsien/goutils/pkgs/ggit"
)

func (that *GSConf) NewGit() *ggit.Git {
	git := ggit.NewGit()
	if os.Getenv(EnableProxyEnvName) != "" {
		git.SetProxyUrl(that.LocalProxy)
	}
	return git
}
