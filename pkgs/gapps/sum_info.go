package gapps

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/gscraper/pkgs/config"
)

const (
	SumType string = "sha256"
)

type ItemInfo struct {
	SHA256  string `json,koanf:"sha256"`
	UpdatAt string `json,koanf:"updat_at"`
}

type InfoResult struct {
	InfoList         map[string]*ItemInfo `json,koanf:"info_list"`
	GVCLatestVersion string               `json,koanf:"gvc_latest_version"`
}

type Info struct {
	Result *InfoResult
	conf   *config.GSConf
	path   string
}

func NewInfo(cnf *config.GSConf) *Info {
	info := &Info{conf: cnf, Result: &InfoResult{InfoList: map[string]*ItemInfo{}}}
	info.path = filepath.Join(cnf.GVCRConifg.GVCResourceDir, config.FileInfoName)
	info.initiate()
	return info
}

func (that *Info) initiate() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		that.Store()
	}
	if ok, _ := utils.PathIsExist(that.path); ok {
		that.Load()
	} else {
		tui.PrintError("Cannot find info files.")
		os.Exit(1)
	}
}

func (that *Info) CheckSum(filename, tempLocalPath string) (updated bool) {
	sumStr := utils.ComputeSum(tempLocalPath, SumType)
	if item := that.Result.InfoList[filename]; item != nil && item.SHA256 == sumStr {
		updated = false
	} else {
		if item == nil {
			item = &ItemInfo{}
			that.Result.InfoList[filename] = item
		}
		item.SHA256 = sumStr
		var cstZone = time.FixedZone("CST", 8*3600)
		item.UpdatAt = time.Now().In(cstZone).Format("2006-01-02 15:04:05")
		updated = true
	}
	return
}

func (that *Info) Store() {
	if content, err := json.MarshalIndent(that.Result, "", "    "); err == nil {
		os.WriteFile(that.path, content, os.ModePerm)
	}
}

func (that *Info) Load() {
	if content, err := os.ReadFile(that.path); err == nil {
		err = json.Unmarshal(content, that.Result)
		if err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}
