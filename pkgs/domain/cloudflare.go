package domain

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	"github.com/moqsien/gscraper/pkgs/conf"
	"github.com/tidwall/gjson"
)

/*
Gather cloudflare third-party domains.
*/

type CFlareDomain struct {
	CNF    *conf.Config
	Result map[string]struct{}
}

func NewCFlareDomain(cnf *conf.Config) (cfd *CFlareDomain) {
	cfd = &CFlareDomain{CNF: cnf}
	cfd.Result = make(map[string]struct{})
	return cfd
}

func (that *CFlareDomain) CheckDomain(domainStr string) {
	conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout:  time.Second * 1,
		Deadline: time.Now().Add(time.Second * 5),
	}, "tcp", fmt.Sprintf("%s:443", domainStr), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(err, " ", domainStr)
		return
	}
	defer conn.Close()
	stats := conn.ConnectionState()
	if certs := stats.PeerCertificates; len(certs) > 0 {
		if certInfo := certs[0]; certInfo != nil {
			s := strings.ToLower(certInfo.Issuer.String())
			if strings.Contains(s, "cloudflare") {
				that.Result[domainStr] = struct{}{}
			}
		}
	}
}

/*
https://www.cloudflare.com/zh-cn/case-studies
*/
func (that *CFlareDomain) GetFromOfficial() (homePages []string) {
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
		fmt.Println("[*] Get from cloudflare: ", slug)
		r := gjson.Get(str, "result.data.caseStudy.homepageUrl")
		homePages = append(homePages, r.String())
	}
	return
}

func (that *CFlareDomain) push() {
	cmdExe := "git"
	if runtime.GOOS == "windows" {
		cmdExe = "git.exe"
	}
	cmd := exec.Command(cmdExe, "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Dir = that.CNF.GvcResourceDir
	err := cmd.Run()
	if err != nil {
		tui.PrintError(err)
		return
	}

	cmd = exec.Command(cmdExe, "commit", "-m", "update")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Dir = that.CNF.GvcResourceDir
	err = cmd.Run()
	if err != nil {
		tui.PrintError(err)
		return
	}

	cmd = exec.Command(cmdExe, "push")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Dir = that.CNF.GvcResourceDir
	err = cmd.Run()
	if err != nil {
		tui.PrintError(err)
		return
	}
}

func (that *CFlareDomain) Run() {
	for _, dStr := range that.CNF.CloudflareDomains {
		that.CheckDomain(dStr)
	}
	urlList := that.GetFromOfficial()
	urlFilePath := filepath.Join(that.CNF.GvcResourceDir, "cloudflare_urls_list.txt")
	os.WriteFile(urlFilePath, []byte(strings.Join(urlList, "\n")), os.ModePerm)
	for _, rUrl := range urlList {
		if u, err := url.Parse(rUrl); err == nil {
			that.CheckDomain(u.Host)
		}
	}
	result := []string{}
	for domain := range that.Result {
		result = append(result, domain)
	}
	content := strings.Join(result, "\n")
	resultFilePath := filepath.Join(that.CNF.GvcResourceDir, that.CNF.CloudflareDomainFile)
	os.WriteFile(resultFilePath, []byte(content), os.ModePerm)
	that.push()
}
