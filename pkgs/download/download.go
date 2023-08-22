package download

import (
	"fmt"
	"os"
	"os/exec"
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

	that.fetcher.SetThreadNum(4)

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

func (that *Downloader) gitPush() {
	cmdName := "git"
	if runtime.GOOS == "windows" {
		cmdName = "git.exe"
	}
	cmd := exec.Command(cmdName, "add", ".")
	cmd.Dir = that.conf.GvcResourceDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		tui.PrintError(err)
		os.Exit(1)
	}

	cmd = exec.Command(cmdName, "commit", "-m", `update`)
	cmd.Dir = that.conf.GvcResourceDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		tui.PrintError(err)
		os.Exit(1)
	}

	cmd = exec.Command(cmdName, "push")
	cmd.Dir = that.conf.GvcResourceDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		tui.PrintError(err)
		os.Exit(1)
	}
}

func (that *Downloader) Start(filenames ...string) {
	if len(filenames) == 0 {
		filenames = that.conf.UrlOrder
	}
	for _, filename := range filenames {
		if dUrl, ok := that.conf.UrlList[filename]; ok {
			that.download(filename, dUrl)
		}
	}
	for _, filename := range filenames {
		if dUrl, ok := that.conf.UrlList[filename]; ok && strings.Contains(dUrl, "gvc_") {
			if tag := that.getLatestTag(dUrl); tag != "" {
				that.info.GVCLatestVersion = tag
				that.info.Store()
			}
			break
		}
	}
	if err := os.RemoveAll(that.getTempDir()); err == nil {
		that.gitPush()
	}
}
