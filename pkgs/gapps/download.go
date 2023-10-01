package gapps

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	archiver "github.com/mholt/archiver/v3"
	"github.com/moqsien/goutils/pkgs/ggit"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/goutils/pkgs/request"
	"github.com/moqsien/gscraper/pkgs/config"
)

/*
download and push files for gvc
*/

var GLOBAL_TO_EXIST bool
var WaitToSweepSig = make(chan struct{})

func IsFileCompressebBySuffix(filename string) bool {
	suffixes := []string{
		".tar", ".gz", ".zip", ".xz", ".7z", ".rar",
	}
	for _, suffix := range suffixes {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}
	return false
}

func GetGVCLocalProjectPath() string {
	home, _ := os.UserHomeDir()
	gvcDir := filepath.Join(home, "data", "projects", "go", "src", "gvc")
	if ok, _ := utils.PathIsExist(gvcDir); ok {
		return gvcDir
	}
	return ""
}

type AppDownloader struct {
	conf    *config.GSConf
	info    *Info
	fetcher *request.Fetcher
	git     *ggit.Git
}

func NewDownloader() (d *AppDownloader) {
	d = &AppDownloader{
		conf: config.NewGSConf(),
		git:  ggit.NewGit(),
	}
	d.info = NewInfo(d.conf)
	d.fetcher = request.NewFetcher()
	return d
}

func (that *AppDownloader) getTempDir() string {
	tempDir := filepath.Join(that.conf.GVCRConifg.GVCResourceDir, "temp")
	if ok, _ := utils.PathIsExist(tempDir); !ok {
		os.MkdirAll(tempDir, 0777)
	}
	return tempDir
}

func (that *AppDownloader) RemoveTempDir() error {
	return os.RemoveAll(that.getTempDir())
}

func (that *AppDownloader) getTempFilePath(filename string) (fPaht string) {
	tempDir := that.getTempDir()
	return filepath.Join(tempDir, filename)
}

func (that *AppDownloader) getLatestTag(dUrl string) (tag string) {
	var _url string
	if strings.Contains(dUrl, "/releases") {
		_url = strings.Split(dUrl, "/releases")[0]
	}
	if strings.Contains(dUrl, "/archive") {
		_url = strings.Split(dUrl, "/archive")[0]
	}
	if _url == "" {
		return
	}
	_url = fmt.Sprintf("%s/releases/latest", _url)
	_url = config.PrepareSubscribeUrl(_url)
	that.fetcher.SetUrl(_url)
	that.fetcher.Timeout = time.Minute
	that.fetcher.RetryTimes = 2
	if resp := that.fetcher.Get(); resp != nil {
		rUrl := resp.RawResponse.Request.URL.String()
		sList := strings.Split(rUrl, "/")
		return sList[len(sList)-1]
	}
	return
}

func (that *AppDownloader) FindUrl(dUrl string) (rUrl string) {
	if strings.Contains(dUrl, "/protobuf/") && strings.Contains(dUrl, `%s`) {
		tag := strings.TrimLeft(that.getLatestTag(dUrl), "v")
		rUrl = fmt.Sprintf(dUrl, tag)
	} else {
		rUrl = dUrl
	}
	return
}

func (that *AppDownloader) download(fileName, dUrl string) {
	_url := that.FindUrl(dUrl)
	if _url == "" {
		return
	}

	that.fetcher.Timeout = time.Minute * 20
	that.fetcher.Url = config.PrepareSubscribeUrl(_url)
	gprint.PrintInfo("[>>>] %s", that.fetcher.Url)

	that.fetcher.SetThreadNum(2)

	tarfile := that.getTempFilePath(fileName)

	var size int64
	if !strings.Contains(dUrl, "refs/heads/") {
		size = that.fetcher.GetAndSaveFile(tarfile, true)
	} else {
		size = that.fetcher.GetFile(tarfile, true)
	}
	if size <= 0 {
		tui.PrintError("Download failed: ", dUrl)
		return
	}
	untarfile := that.getTempFilePath(strings.ReplaceAll(fileName, ".", "_"))
	toCopy := true
	if IsFileCompressebBySuffix(untarfile) && archiver.Unarchive(tarfile, untarfile) != nil {
		toCopy = false
	}
	if toCopy {
		os.RemoveAll(untarfile)
		if sumChanged := that.info.CheckSum(fileName, tarfile); sumChanged {
			if _, err := utils.CopyFile(tarfile, filepath.Join(that.conf.GVCRConifg.GVCResourceDir, fileName)); err == nil {
				that.info.Store()
			} else {
				that.info.Load()
			}
		}
	}
}

func (that *AppDownloader) copyGVCFromLocal(fileName string) {
	srcFile := filepath.Join(GetGVCLocalProjectPath(), "build", fileName)
	tui.PrintInfo(fmt.Sprintf("copying: %s", srcFile))
	if sumChanged := that.info.CheckSum(fileName, srcFile); sumChanged {
		if _, err := utils.CopyFile(srcFile, filepath.Join(that.conf.GVCRConifg.GVCResourceDir, fileName)); err == nil {
			that.updateGVCVersionFromLocal()
			that.info.Store()
		} else {
			that.info.Load()
		}
	}
}

func (that *AppDownloader) updateGVCVersionFromLocal() {
	r, _ := utils.ExecuteSysCommand(true, GetGVCLocalProjectPath(), "git", "describe", "--abbrev=0", "--tags")
	content, _ := io.ReadAll(r)
	if len(content) > 0 {
		that.info.Result.GVCLatestVersion = strings.Trim(string(content), "\n")
	}
}

func (that *AppDownloader) GetGitCmd() string {
	cmdName := "git"
	if runtime.GOOS == "windows" {
		cmdName = "git.exe"
	}
	return cmdName
}

func (that *AppDownloader) Start(filenames ...string) {
	that.git.SetWorkDir(that.conf.GVCRConifg.GVCResourceDir)
	that.git.PullBySSH()

	itemList := map[string]string{}
	nameList := []string{}
	maxIndex := config.GetMax(that.conf.GVCRConifg.APPUrls)
	for i := 0; i <= maxIndex; i++ {
		if item := that.conf.GVCRConifg.APPUrls[i]; item.FileName != "" && item.Url != "" {
			itemList[item.FileName] = item.Url
			nameList = append(nameList, item.FileName)
		}
	}
	if len(filenames) == 0 {
		filenames = nameList
	}

	doExists, _ := utils.PathIsExist(GetGVCLocalProjectPath())
	for _, filename := range filenames {
		if GLOBAL_TO_EXIST {
			close(WaitToSweepSig)
			return
		}
		if dUrl, ok := itemList[filename]; ok && !strings.Contains(dUrl, "gvc_") {
			that.download(filename, dUrl)
		} else if strings.Contains(dUrl, "gvc_") && doExists {
			that.copyGVCFromLocal(filename)
		} else if strings.Contains(dUrl, "gvc_") && !doExists {
			that.download(filename, dUrl)
		}
	}
	if !doExists {
		for _, filename := range filenames {
			if GLOBAL_TO_EXIST {
				close(WaitToSweepSig)
				return
			}
			if dUrl, ok := itemList[filename]; ok && strings.Contains(dUrl, "gvc_") {
				if tag := that.getLatestTag(dUrl); tag != "" {
					that.info.Result.GVCLatestVersion = tag
					that.info.Store()
				}
				break
			}
		}
	} else {
		output, _ := utils.ExecuteSysCommand(true, GetGVCLocalProjectPath(), that.GetGitCmd(), "describe", "--abbrev=0", "--tags")
		that.info.Result.GVCLatestVersion = strings.TrimRight(output.String(), "\n")
	}

	if err := os.RemoveAll(that.getTempDir()); err == nil {
		that.Push()
	}
}

func (that *AppDownloader) Push() {
	that.git.CommitAndPush("update softwares")
}
