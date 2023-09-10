package download

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	archiver "github.com/mholt/archiver/v3"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/goutils/pkgs/request"
	"github.com/moqsien/gscraper/pkgs/conf"
)

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

type Downloader struct {
	conf    *conf.Config
	info    *Info
	fetcher *request.Fetcher
}

func NewDownloader() (d *Downloader) {
	d = &Downloader{
		conf: conf.NewConfig(),
	}
	d.info = NewInfo(d.conf)
	d.fetcher = request.NewFetcher()
	return d
}

func (that *Downloader) getTempDir() string {
	tempDir := filepath.Join(that.conf.GvcResourceDir, "temp")
	if ok, _ := utils.PathIsExist(tempDir); !ok {
		os.MkdirAll(tempDir, 0777)
	}
	return tempDir
}

func (that *Downloader) RemoveTempDir() error {
	return os.RemoveAll(that.getTempDir())
}

func (that *Downloader) getTempFilePath(filename string) (fPaht string) {
	tempDir := that.getTempDir()
	return filepath.Join(tempDir, filename)
}

func (that *Downloader) getLatestTag(dUrl string) (tag string) {
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
	that.fetcher.SetUrl(that.conf.GithubSpeedupUrl + _url)
	that.fetcher.Timeout = time.Minute
	that.fetcher.RetryTimes = 2
	if resp := that.fetcher.Get(); resp != nil {
		rUrl := resp.RawResponse.Request.URL.String()
		sList := strings.Split(rUrl, "/")
		return sList[len(sList)-1]
	}
	return
}

func (that *Downloader) FindUrl(dUrl string) (rUrl string) {
	if strings.Contains(dUrl, "/protobuf/") && strings.Contains(dUrl, `%s`) {
		tag := strings.TrimLeft(that.getLatestTag(dUrl), "v")
		rUrl = fmt.Sprintf(dUrl, tag)
	} else {
		rUrl = dUrl
	}
	return
}

func (that *Downloader) download(fileName, dUrl string) {
	_url := that.FindUrl(dUrl)
	if _url == "" {
		return
	}

	that.fetcher.Timeout = time.Minute * 20
	that.fetcher.Url = that.conf.GithubSpeedupUrl + _url
	tui.PrintInfo("[>>>] ", that.fetcher.Url)

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
			if _, err := utils.CopyFile(tarfile, filepath.Join(that.conf.GvcResourceDir, fileName)); err == nil {
				that.info.Store()
			} else {
				that.info.Load()
			}
		}
	}
}

func (that *Downloader) copyGVCFromLocal(fileName string) {
	srcFile := filepath.Join(GetGVCLocalProjectPath(), "build", fileName)
	tui.PrintInfo(fmt.Sprintf("copying: %s", srcFile))
	if sumChanged := that.info.CheckSum(fileName, srcFile); sumChanged {
		if _, err := utils.CopyFile(srcFile, filepath.Join(that.conf.GvcResourceDir, fileName)); err == nil {
			that.updateGVCVersionFromLocal()
			that.info.Store()
		} else {
			that.info.Load()
		}
	}
}

func (that *Downloader) updateGVCVersionFromLocal() {
	r, _ := utils.ExecuteSysCommand(true, GetGVCLocalProjectPath(), "git", "describe", "--abbrev=0", "--tags")
	content, _ := io.ReadAll(r)
	if len(content) > 0 {
		that.info.GVCLatestVersion = string(content)
	}
}

func (that *Downloader) GetGitCmd() string {
	cmdName := "git"
	if runtime.GOOS == "windows" {
		cmdName = "git.exe"
	}
	return cmdName
}

func (that *Downloader) GitPush() {
	cmdName := that.GetGitCmd()
	_, err := utils.ExecuteSysCommand(false, that.conf.GvcResourceDir, cmdName, "add", ".")
	if err != nil {
		tui.PrintError(err)
		os.Exit(1)
	}

	_, err = utils.ExecuteSysCommand(false, that.conf.GvcResourceDir, cmdName, "commit", "-m", `update`)
	if err != nil {
		tui.PrintError(err)
		os.Exit(1)
	}

	_, err = utils.ExecuteSysCommand(false, that.conf.GvcResourceDir, cmdName, "push")
	if err != nil {
		tui.PrintError(err)
		os.Exit(1)
	}
}

func (that *Downloader) Start(filenames ...string) {
	if len(filenames) == 0 {
		filenames = that.conf.UrlOrder
	}

	doExists, _ := utils.PathIsExist(GetGVCLocalProjectPath())
	for _, filename := range filenames {
		if GLOBAL_TO_EXIST {
			close(WaitToSweepSig)
			return
		}
		if dUrl, ok := that.conf.UrlList[filename]; ok && !strings.Contains(dUrl, "gvc_") {
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
			if dUrl, ok := that.conf.UrlList[filename]; ok && strings.Contains(dUrl, "gvc_") {
				if tag := that.getLatestTag(dUrl); tag != "" {
					that.info.GVCLatestVersion = tag
					that.info.Store()
				}
				break
			}
		}
	} else {
		output, _ := utils.ExecuteSysCommand(true, GetGVCLocalProjectPath(), that.GetGitCmd(), "describe", "--abbrev=0", "--tags")
		that.info.GVCLatestVersion = strings.TrimRight(output.String(), "\n")
	}

	if err := os.RemoveAll(that.getTempDir()); err == nil {
		that.GitPush()
	}
}
