package utils

import (
	"fmt"
	"os"
)

func GetVersion() string {
	// 读取版本文件
	versionBytes, err := os.ReadFile("VERSION")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 将版本号转换为字符串
	version := string(versionBytes)

	return version
}
