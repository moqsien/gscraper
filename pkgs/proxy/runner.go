package proxy

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/moqsien/goutils/pkgs/crypt"
	"github.com/moqsien/goutils/pkgs/ggit"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gscraper/pkgs/config"
	"github.com/moqsien/vpnparser/pkgs/outbound"
)

func HandleQuery(rawUri string) (result string) {
	result = rawUri
	if !strings.Contains(rawUri, "?") {
		return
	}
	sList := strings.Split(rawUri, "?")
	query := sList[1]
	if strings.Contains(query, ";") && !strings.Contains(query, "&") {
		result = sList[0] + "?" + strings.ReplaceAll(sList[1], ";", "&")
	}
	return
}

func ParseRawUri(rawUri string) (result string) {
	if strings.HasPrefix(rawUri, "vmess://") {
		if r := crypt.DecodeBase64(strings.Split(rawUri, "://")[1]); r != "" {
			result = "vmess://" + r
		}
		return
	}

	if strings.Contains(rawUri, "\u0026") {
		rawUri = strings.ReplaceAll(rawUri, "\u0026", "&")
	}
	if strings.Contains(rawUri, "amp;") {
		rawUri = strings.ReplaceAll(rawUri, "amp;", "")
	}
	rawUri, _ = url.QueryUnescape(rawUri)
	r, err := url.Parse(rawUri)
	result = rawUri
	if err != nil {
		gprint.PrintError("%+v", err)
		return
	}

	host := r.Host
	uname := r.User.Username()
	passw, hasPassword := r.User.Password()

	if !strings.Contains(rawUri, "@") {
		if hostDecrypted := crypt.DecodeBase64(host); hostDecrypted != "" {
			result = strings.ReplaceAll(rawUri, host, hostDecrypted)
		}
	} else if uname != "" && !hasPassword && !strings.Contains(uname, "-") {
		if unameDecrypted := crypt.DecodeBase64(uname); unameDecrypted != "" {
			result = strings.ReplaceAll(rawUri, uname, unameDecrypted)
		}
	} else {
		if passwDecrypted := crypt.DecodeBase64(passw); passwDecrypted != "" {
			result = strings.ReplaceAll(rawUri, passw, passwDecrypted)
		}
	}

	if strings.Contains(result, "%") {
		result, _ = url.QueryUnescape(result)
	}
	result = HandleQuery(result)
	if strings.Contains(result, "127.0.0.1") || strings.Contains(result, "127.0.0.0") {
		return ""
	}
	return
}

type ISite interface {
	SetHandler(handler func([]string))
	Run()
	Type() string
}

type ProxyRunner struct {
	Result *outbound.Result `json:"vpn_list"`
	cnf    *config.GSConf
	sites  []ISite
	r      map[string]struct{}
	git    *ggit.Git
}

func NewProxyRunner() (pr *ProxyRunner) {
	pr = &ProxyRunner{
		Result: outbound.NewResult(),
		cnf:    config.NewGSConf(),
		git:    ggit.NewGit(),
	}
	return
}

func (that *ProxyRunner) AddSite(site ISite) {
	that.sites = append(that.sites, site)
}

func (that *ProxyRunner) wrapItem(rawUri string) *outbound.ProxyItem {
	item := outbound.NewItem(rawUri)
	if strings.HasPrefix(item.Address, "127.0.") {
		return nil
	}
	item.GetOutbound()
	return item
}

func (that *ProxyRunner) Run() {
	doProxy := false
	that.r = map[string]struct{}{}
	that.Result = outbound.NewResult()
	that.git.SetWorkDir(that.cnf.NeoboxRConfig.NeoboxResourceDir)
	that.git.PullBySSH()
	for _, site := range that.sites {
		switch site.Type() {
		case "proxies":
			site.SetHandler(func(result []string) {
				for _, rawUri := range result {
					rawUri = ParseRawUri(rawUri)
					proxyItem := that.wrapItem(rawUri)
					proxyStr := fmt.Sprintf("%s%s:%d", proxyItem.Scheme, proxyItem.Address, proxyItem.Port)
					if _, ok := that.r[proxyStr]; !ok {
						that.Result.AddItem(proxyItem)
						that.r[proxyStr] = struct{}{}
					}
				}
			})
			doProxy = true
		case "domains":
			site.SetHandler(func(result []string) {
				gprint.PrintInfo("Find %d available domains", len(result))
			})
		default:
		}
		site.Run()
	}
	if doProxy {
		gprint.PrintSuccess("Total Proxies: %d", that.Result.Len())
		gprint.PrintSuccess(
			"vmess[%d]; vless[%d]; ss[%d]; trojan[%d]; ssr[%d]",
			that.Result.VmessTotal,
			that.Result.VlessTotal,
			that.Result.SSTotal,
			that.Result.TrojanTotal,
			that.Result.SSRTotal,
		)
		fPath := filepath.Join(that.cnf.NeoboxRConfig.NeoboxResourceDir, config.NeoboxResultFileName)
		var cstZone = time.FixedZone("CST", 8*3600)
		now := time.Now().In(cstZone)
		that.Result.UpdateAt = now.Format("2006-01-02 15:04:05")
		if that.Result.Len() <= 0 {
			return
		}
		if content, err := json.Marshal(that.Result); err == nil {
			gprint.PrintWarning("neobox key: %s", that.cnf.NeoboxRConfig.NeoboxKey)
			cc := crypt.NewCrptWithKey([]byte(that.cnf.NeoboxRConfig.NeoboxKey))
			if r, err := cc.AesEncrypt([]byte(content)); err == nil {
				os.WriteFile(fPath, r, os.ModePerm)
			}
		}
	}
	gprint.PrintInfo("push to remote repository...")
	that.git.CommitAndPush("update")
}
