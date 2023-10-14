package utils

import (
	"net/url"
	"strings"
)

func getUrlFromParam(urlStr string) (string, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	fromParam := parsed.Query().Get("from")
	return strings.TrimSpace(fromParam), nil
}

func getUrlOssParam(urlStr string) (string, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	fromParam := parsed.Query().Get("oss")
	return strings.TrimSpace(fromParam), nil
}
