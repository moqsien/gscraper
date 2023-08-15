package main

import (
	"github.com/moqsien/gscraper/pkgs/download"
)

func main() {
	d := download.NewDownloader()
	d.Run()
}
