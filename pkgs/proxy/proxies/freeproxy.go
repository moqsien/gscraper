package proxies

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/moqsien/goutils/pkgs/crypt"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	"github.com/moqsien/gscraper/pkgs/config"
)

// https://wanshanziwo.eu.org/
type WSZiwo struct {
	result  []string
	fetcher *request.Fetcher
	handler func([]string)
	CNF     *config.GSConf
	url     string
	randStr string
}

func NewWSZiwo() (wsz *WSZiwo) {
	wsz = &WSZiwo{
		fetcher: request.NewFetcher(),
		CNF:     config.NewGSConf(),
		url:     "https://wanshanziwo.eu.org/",
	}
	return
}

func (that *WSZiwo) SetHandler(handler func([]string)) {
	that.handler = handler
}

func (that *WSZiwo) getRand() {
	that.fetcher.SetUrl(that.url)
	content, _ := that.fetcher.GetString()
	if doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(content)); err == nil {
		doc.Find("table").Find("td").First().Each(func(i int, s *goquery.Selection) {
			sUrl := s.Text()
			if that.randStr == "" && strings.HasPrefix(sUrl, "https://") {
				if u, err := url.Parse(sUrl); err == nil {
					that.randStr = u.Query().Get("rand")
				}
			}
		})
	}
}

func (that *WSZiwo) getWSZiwo() {
	list := []string{
		"ss/sub?rand=%s",
		"vmess/sub?rand=%s",
		"airport/sub?rand=%s",
		"airport?rand=%s&core=clash",
	}
	that.getRand()
	if that.randStr != "" {
		for _, q := range list {
			sUrl := that.url + fmt.Sprintf(q, that.randStr)
			gprint.PrintInfo("Download: %s", sUrl)

			that.fetcher.Headers = map[string]string{
				"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.43",
			}
			urlPattern := `https://api.dler.io/sub?target=v2ray&insert=false&config=https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/config/ACL4SSR_Online_NoAuto.ini&url={url}`
			if !strings.Contains(sUrl, "/ss/sub?") {
				sUrl = strings.ReplaceAll(urlPattern, "{url}", sUrl)
			}
			that.fetcher.SetUrl(sUrl)
			if content, _ := that.fetcher.GetString(); len(content) > 0 {
				// fmt.Println(content)
				decryptedContent := crypt.DecodeBase64(content)
				for _, item := range strings.Split(decryptedContent, "\n") {
					if strings.Contains(item, "://") {
						that.result = append(that.result, item)
					}
				}
			}
		}
	}
}

func (that *WSZiwo) Run() {
	that.getWSZiwo()
	if that.handler != nil {
		that.handler(that.result)
	}
}
