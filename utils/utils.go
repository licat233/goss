package utils

import (
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// 检查是否为URL
func IsURL(input string) bool {
	u, err := url.Parse(input)
	if err != nil {
		return false
	}
	//存在协议，一定是url
	if u.Scheme != "" {
		return true
	}
	//没有Host，一定不是url
	if u.Host == "" {
		return false
	}
	//否则就一定是url

	return true
}

func UUID() string {
	return uuid.New().String()
}

func FileExt(filename string) string {
	return filepath.Ext(filename)
}

func UUIDhex() string {
	u := uuid.New().String()
	uidWithoutDash := strings.ReplaceAll(u, "-", "")
	return uidWithoutDash
}
