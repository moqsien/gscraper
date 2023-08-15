package download

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	tui "github.com/moqsien/goutils/pkgs/gtui"
)

func CheckSum(fpath, cType, cSum string) (r bool) {
	if cSum != ComputeSum(fpath, cType) {
		tui.PrintError("Checksum failed.")
		return
	}
	tui.PrintSuccess("Checksum succeeded.")
	return true
}

func ComputeSum(fpath, sumType string) (sumStr string) {
	f, err := os.Open(fpath)
	if err != nil {
		tui.PrintError(fmt.Sprintf("Open file failed: %+v", err))
		return
	}
	defer f.Close()

	var h hash.Hash
	switch strings.ToLower(sumType) {
	case "sha256":
		h = sha256.New()
	case "sha1":
		h = sha1.New()
	case "sha512":
		h = sha512.New()
	default:
		tui.PrintError(fmt.Sprintf("[Crypto] %s is not supported.", sumType))
		return
	}

	if _, err = io.Copy(h, f); err != nil {
		tui.PrintError(fmt.Sprintf("Copy file failed: %+v", err))
		return
	}

	sumStr = hex.EncodeToString(h.Sum(nil))
	return
}
