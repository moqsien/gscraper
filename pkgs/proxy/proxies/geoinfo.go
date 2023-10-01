package proxies

import (
	"path/filepath"
	"time"

	"github.com/moqsien/goutils/pkgs/request"
	"github.com/moqsien/gscraper/pkgs/config"
)

type GeoInfo struct {
	fetcher *request.Fetcher
	CNF     *config.GSConf
	handler func([]string)
}

func NewGeoInfo() (gi *GeoInfo) {
	gi = &GeoInfo{
		fetcher: request.NewFetcher(),
		CNF:     config.NewGSConf(),
	}
	return
}

func (that *GeoInfo) SetHandler(handler func([]string)) {
	that.handler = handler
}

func (that *GeoInfo) download() {
	maxIndex := config.GetMax(that.CNF.NeoboxRConfig.GeoInfoUrls)
	for i := 0; i <= maxIndex; i++ {
		if item := that.CNF.NeoboxRConfig.GeoInfoUrls[i]; item.FileName != "" && item.Url != "" {
			fPath := filepath.Join(that.CNF.NeoboxRConfig.NeoboxResourceDir, item.FileName)
			iUrl := config.PrepareSubscribeUrl(item.Url)
			that.fetcher.SetUrl(iUrl)
			that.fetcher.SetThreadNum(2)
			that.fetcher.Timeout = 15 * time.Minute
			that.fetcher.GetAndSaveFile(fPath, true)
		}
	}
}

func (that *GeoInfo) Run() {
	that.download()
	if that.handler != nil {
		that.handler([]string{})
	}
}
