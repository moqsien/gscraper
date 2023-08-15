package download

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	archiver "github.com/mholt/archiver/v3"
	tui "github.com/moqsien/goutils/pkgs/gtui"
	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/goutils/pkgs/request"
	"github.com/moqsien/gscraper/pkgs/conf"
)

func CopyFile(src, dst string) (written int64, err error) {
	srcFile, err := os.Open(src)

	if err != nil {
		tui.PrintError(fmt.Sprintf("Cannot open file: %+v", err))
		return
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		tui.PrintError(fmt.Sprintf("Cannot open file: %+v", err))
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
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

func (that *Downloader) download(fileName, dUrl string) {
	that.fetcher.Timeout = time.Minute * 20
	that.fetcher.Url = that.conf.GithubSpeedupUrl + dUrl
	tui.PrintInfo(that.fetcher.Url)
	if !strings.Contains(that.fetcher.Url, "master") {
		return
	}
	that.fetcher.SetThreadNum(4)

	tarfile := that.getTempFilePath(fileName)
	if size := that.fetcher.GetAndSaveFile(tarfile, true); size > 0 {
		untarfile := that.getTempFilePath(strings.ReplaceAll(fileName, ".", "_"))
		if err := archiver.Unarchive(tarfile, untarfile); err == nil {
			os.RemoveAll(untarfile)
			if sumChanged := that.info.CheckSum(fileName, tarfile); sumChanged {
				_, err = CopyFile(tarfile, filepath.Join(that.conf.GvcResourceDir, fileName))
				if err == nil {
					that.info.Store()
				} else {
					that.info.Load()
				}
			}
		} else {
			that.info.Load()
		}
	}
}

func (that *Downloader) Run() {
	for filename, dUrl := range that.conf.UrlList {
		that.download(filename, dUrl)
	}
	os.RemoveAll(that.getTempDir())
}
