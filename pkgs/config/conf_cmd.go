package config

import (
	"fmt"
	"strings"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gutils"
)

/*
edit configfile by cmds
*/

func getNameFromUrl(dUrl string) string {
	sList := strings.Split(dUrl, "/")
	filename := sList[len(sList)-1]
	if strings.Contains(filename, "master") || strings.Contains(filename, "main") {
		sList = strings.Split(dUrl, "github.com/")
		if len(sList) == 2 {
			s := sList[1]
			sList = strings.Split(s, "/")
			if len(sList) > 1 {
				return fmt.Sprintf("%s.zip", sList[1])
			}
		}
	}
	return filename
}

type MapValue interface {
	string | FileUrl
}

func getMax[T MapValue](list map[int]T) (maxIdx int) {
	for idx := range list {
		if idx > maxIdx {
			maxIdx = idx
		}
	}
	return
}

func (that *GSConf) AddSubscribeUrlForNeobox(subUrl string) {
	that.Load()
	idx := getMax(that.NeoboxRConfig.ProxySubcribeUrlList) + 1
	that.NeoboxRConfig.ProxySubcribeUrlList[idx] = subUrl
	that.Save()
}

func (that *GSConf) DelSubscribeUrlForNeobox(subUrl string) {
	that.Load()
	for idx, item := range that.NeoboxRConfig.ProxySubcribeUrlList {
		if item == subUrl {
			delete(that.NeoboxRConfig.ProxySubcribeUrlList, idx)
		}
	}
	that.Save()
}

func (that *GSConf) AddGVCAppUrl(appUrl string) {
	that.Load()
	idx := getMax(that.GVCRConifg.APPUrls) + 1
	fileName := getNameFromUrl(appUrl)
	that.GVCRConifg.APPUrls[idx] = FileUrl{Url: appUrl, FileName: fileName}
	that.Save()
}

func (that *GSConf) DelGVCAppUrl(appNameOrUrl string) {
	that.Load()
	for idx, f := range that.GVCRConifg.APPUrls {
		if f.FileName == appNameOrUrl || f.Url == appNameOrUrl {
			delete(that.GVCRConifg.APPUrls, idx)
		}
	}
	that.Save()
}

func (that *GSConf) ShowAppUrls() {
	that.Load()
	maxIndex := getMax(that.GVCRConifg.APPUrls)
	for i := 0; i < maxIndex; i++ {
		if item := that.GVCRConifg.APPUrls[i]; item.Url != "" && item.FileName != "" {
			gprint.PrintInfo("%s: %s", item.FileName, item.Url)
		}
	}
}

func (that *GSConf) SetNeoboxKey() {
	that.NeoboxRConfig.NeoboxKey = gutils.RandomString(16)
	that.ShowNeoboxKey()
}

func (that *GSConf) ShowNeoboxKey() {
	gprint.PrintInfo(that.NeoboxRConfig.NeoboxKey)
}
