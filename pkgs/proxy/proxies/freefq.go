package proxies

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	"github.com/moqsien/gscraper/pkgs/config"
)

type FreeFQ struct {
	result  []string
	fetcher *request.Fetcher
	handler func([]string)
	CNF     *config.GSConf
	urls    []string
}

func NewFreeFQ() (ffq *FreeFQ) {
	ffq = &FreeFQ{
		fetcher: request.NewFetcher(),
		CNF:     config.NewGSConf(),
		urls: []string{
			"https://freefq.com/v2ray/",
			"https://freefq.com/free-xray/",
			"https://freefq.com/free-ss/",
			"https://freefq.com/free-trojan/",
			"https://freefq.com/free-ssr/",
		},
	}

	proxyEnvName := "FREE_FQ_PROXY"
	os.Setenv(proxyEnvName, "http://127.0.0.1:2023")
	ffq.fetcher.SetProxyEnvName(proxyEnvName)
	return
}

func (that *FreeFQ) SetHandler(handler func([]string)) {
	that.handler = handler
}

func (that *FreeFQ) getUrl(sUrl string) (r string) {
	var c *http.Client
	if os.Getenv(config.EnableProxyEnvName) != "" {
		u, _ := url.Parse("htttp://127.0.0.1:2023")
		t := &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 10,
			IdleConnTimeout: time.Duration(10) * time.Second,
			Proxy:           http.ProxyURL(u),
		}
		c = &http.Client{
			Transport: t,
			Timeout:   time.Duration(30) * time.Second,
		}
	} else {
		c = &http.Client{
			Timeout: time.Duration(30) * time.Second,
		}
	}

	if resp, err := c.Get(sUrl); err == nil {
		content, _ := io.ReadAll(resp.Body)
		r = string(content)
		resp.Body.Close()
	} else {
		fmt.Println(err)
	}
	return
}

func (that *FreeFQ) getUrls() (r []string) {
	for _, sUrl := range that.urls {

		content := that.getUrl(sUrl)
		// fmt.Println(content)
		if doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(content)); err == nil && doc != nil {
			href := doc.Find("td.news_list").Find("ul").First().Find("li").First().Find("a").AttrOr("href", "")
			if href != "" {
				detailUrl := "https://freefq.com" + href
				gprint.PrintInfo("Dowload: %s", detailUrl)
				content = that.getUrl(detailUrl)
				if doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(content)); err == nil && doc != nil {
					fUrl := doc.Find("fieldset").Find("a").AttrOr("href", "")
					if fUrl != "" {
						r = append(r, fUrl)
					}
				}
			}
		}
	}
	return
}

func (that *FreeFQ) getRawUris() {
	urls := that.getUrls()
	for _, u := range urls {
		content := that.getUrl(u)
		for _, rawUri := range strings.Split(content, "\n") {
			rawUri = strings.TrimSpace(rawUri)
			rawUri = strings.ReplaceAll(rawUri, "<br>", "")
			rawUri = strings.ReplaceAll(rawUri, "</p>", "")
			rawUri = strings.ReplaceAll(rawUri, "<p>", "")
			if strings.Contains(rawUri, "<script") {
				continue
			}
			if strings.HasPrefix(rawUri, "vmess://") {
				that.result = append(that.result, rawUri)
			} else if strings.HasPrefix(rawUri, "vless://") {
				that.result = append(that.result, rawUri)
			} else if strings.HasPrefix(rawUri, "ss://") {
				that.result = append(that.result, rawUri)
			} else if strings.HasPrefix(rawUri, "ssr://") {
				that.result = append(that.result, rawUri)
			} else if strings.HasPrefix(rawUri, "trojan://") {
				that.result = append(that.result, rawUri)
			}
		}
	}
}

func (that *FreeFQ) Run() {
	that.getRawUris()
	if that.handler != nil {
		that.handler(that.result)
	}
}
