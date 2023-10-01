package proxies

import (
	"fmt"
	"strings"
	"time"

	"github.com/moqsien/goutils/pkgs/crypt"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	"github.com/moqsien/gscraper/pkgs/config"
)

type Subscribers struct {
	result  []string
	fetcher *request.Fetcher
	handler func([]string)
	CNF     *config.GSConf
}

func NewSubscribers() (sub *Subscribers) {
	sub = &Subscribers{
		fetcher: request.NewFetcher(),
		CNF:     config.NewGSConf(),
	}
	sub.fetcher.Timeout = 30 * time.Second
	return
}

func (that *Subscribers) SetHandler(handler func([]string)) {
	that.handler = handler
}

func (that *Subscribers) Type() string {
	return "proxies"
}

func (that *Subscribers) getSubUrl() {
	maxIndex := config.GetMax(that.CNF.NeoboxRConfig.ProxySubcribeUrlList)
	for i := 0; i <= maxIndex; i++ {
		if item := that.CNF.NeoboxRConfig.ProxySubcribeUrlList[i]; item != "" {
			sUrl := config.PrepareSubscribeUrl(item)
			gprint.PrintInfo("Download: %s", sUrl)
			that.fetcher.SetUrl(sUrl)
			if content, statusCode := that.fetcher.GetString(); len(content) > 0 {
				decryptedContent := crypt.DecodeBase64(content)
				if len(decryptedContent) == 0 && len(content) > 500 && !strings.Contains(content, "</html>") {
					fmt.Println(content)
					for _, encryptedContent := range strings.Split(content, "\n") {
						decryptedContent = crypt.DecodeBase64(strings.TrimSpace(encryptedContent))
						for _, rawUri := range strings.Split(decryptedContent, "\n") {
							if strings.Contains(rawUri, "://") {
								that.result = append(that.result, strings.TrimSpace(rawUri))
							}
						}
					}
				} else if len(content) > 800 && !strings.Contains(content, "</html>") {
					for _, rawUri := range strings.Split(decryptedContent, "\n") {
						if strings.Contains(rawUri, "://") {
							that.result = append(that.result, strings.TrimSpace(rawUri))
						}
					}
				}
			} else {
				gprint.PrintError("status code: %d", statusCode)
			}
		}
	}
}

func (that *Subscribers) Run() {
	that.getSubUrl()
	if that.handler != nil {
		that.handler(that.result)
	}
}
