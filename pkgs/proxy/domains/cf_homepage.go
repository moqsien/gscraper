package domains

import (
	"fmt"
	"net/url"
	"time"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/request"
	"github.com/tidwall/gjson"
)

/*
https://www.cloudflare.com/zh-cn/case-studies
*/
func GetFromOfficial() (hostList []string) {
	f := request.NewFetcher()
	f.Timeout = 60 * time.Second
	f.SetUrl("https://www.cloudflare.com/page-data/zh-cn/case-studies/page-data.json")
	str, _ := f.GetString()
	if str == "" {
		return
	}
	rList := gjson.Get(str, "result.data.caseStudies.nodes.#.nameUrlSlug").Array()
	for _, result := range rList {
		slug := result.String()
		if slug == "" {
			return
		}
		pageUrl := fmt.Sprintf("https://www.cloudflare.com/page-data/zh-cn/case-studies/%s/page-data.json", slug)
		f.SetUrl(pageUrl)
		str, _ = f.GetString()
		if str == "" {
			return
		}
		// fmt.Println("[*] Get from cloudflare: ", slug)
		gprint.PrintInfo("Get domains from cloudflare: %s", slug)
		r := gjson.Get(str, "result.data.caseStudy.homepageUrl")
		if u, err := url.Parse(r.String()); err == nil && u != nil && u.Host != "" {
			hostList = append(hostList, u.Host)
		}
	}
	return
}
