package domains

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gscraper/pkgs/config"
)

type CFDomains struct {
	CNF      *config.GSConf
	Result   map[string]struct{}
	handler  func([]string)
	lock     *sync.Mutex
	sendChan chan string
}

func NewCFDomains() (cfd *CFDomains) {
	cfd = &CFDomains{
		CNF:      config.NewGSConf(),
		Result:   map[string]struct{}{},
		lock:     &sync.Mutex{},
		sendChan: make(chan string, 10),
	}
	return cfd
}

func (that *CFDomains) send(hostList []string) {
	that.sendChan = make(chan string, 10)
	for _, host := range hostList {
		that.sendChan <- host
	}
	close(that.sendChan)
}

func (that *CFDomains) SetHandler(handler func([]string)) {
	that.handler = handler
}

func (that *CFDomains) Type() string {
	return "domains"
}

func (that *CFDomains) CheckDomain(domainStr string) {
	conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout:  time.Second * 2,
		Deadline: time.Now().Add(time.Second * 6),
	}, "tcp", fmt.Sprintf("%s:443", domainStr), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		gprint.PrintError("Domain: %s, err: %+v", domainStr, err)
		return
	}
	defer conn.Close()
	stats := conn.ConnectionState()
	if certs := stats.PeerCertificates; len(certs) > 0 {
		if certInfo := certs[0]; certInfo != nil {
			s := strings.ToLower(certInfo.Issuer.String())
			if strings.Contains(s, "cloudflare") {
				that.lock.Lock()
				that.Result[domainStr] = struct{}{}
				that.lock.Unlock()
			}
		}
	}
}

func (that *CFDomains) handleDomains() {
	hostList := GetFromOfficial()
	hostList2 := strings.Split(config.CF_Raw_Domains, "\n")
	hostList = append(hostList, hostList2...)

	os.WriteFile(filepath.Join(
		that.CNF.NeoboxRConfig.NeoboxResourceDir, config.CloudflareRawDomainFileName),
		[]byte(strings.Join(hostList, "\n")),
		os.ModePerm,
	)
	go that.send(hostList)
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					gprint.PrintError("%+v", r)
				}
				wg.Done()
			}()
			for {
				select {
				case domainStr, ok := <-that.sendChan:
					if !ok {
						return
					}
					that.CheckDomain(domainStr)
				default:
					time.Sleep(time.Millisecond * 100)
				}
			}
		}()
	}
	wg.Wait()

	rList := []string{}
	for domain := range that.Result {
		rList = append(rList, domain)
	}

	r := strings.Join(rList, "\n")
	fPath := filepath.Join(that.CNF.NeoboxRConfig.NeoboxResourceDir, config.CloudflareDomainFileName)
	os.WriteFile(fPath, []byte(r), os.ModePerm)
}

func (that *CFDomains) Run() {
	that.handleDomains()
	result := []string{}
	for domain := range that.Result {
		result = append(result, domain)
	}
	if that.handler != nil {
		that.handler(result)
	}
}
