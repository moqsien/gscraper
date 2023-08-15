package download

import (
	"os"
	"path/filepath"
	"time"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/moqsien/gscraper/pkgs/conf"
)

const (
	SumType string = "sha256"
)

type ItemInfo struct {
	SHA256  string `json,koanf:"sha256"`
	UpdatAt string `json,koanf:"updat_at"`
}

type Info struct {
	InfoList map[string]*ItemInfo `json,koanf:"info_list"`
	koanfer  *koanfer.JsonKoanfer
	conf     *conf.Config
	path     string
}

func NewInfo(cnf *conf.Config) *Info {
	info := &Info{conf: cnf, InfoList: map[string]*ItemInfo{}}
	info.path = filepath.Join(cnf.GvcResourceDir, "files_info.json")
	info.koanfer, _ = koanfer.NewKoanfer(info.path)
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
	sumStr := ComputeSum(tempLocalPath, SumType)
	if item := that.InfoList[filename]; item != nil && item.SHA256 == sumStr {
		updated = false
	} else {
		if item == nil {
			item = &ItemInfo{}
			that.InfoList[filename] = item
		}
		item.SHA256 = sumStr
		item.UpdatAt = time.Now().Format("2006-01-02 15:04:05")
		updated = true
	}
	return
}

func (that *Info) Store() {
	that.koanfer.Save(that)
}

func (that *Info) Load() {
	that.koanfer.Load(that)
}
