package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// 获取指定目录下的某些格式的文件列表
func GetDirFiles(dir string, ext string) ([]string, error) {
	if dir == "" {
		dir = "."
	}
	if ext != "" {
		if strings.HasPrefix(ext, ".") {
			ext = strings.TrimLeft(ext, ".")
		}
		ext = "." + ext
	}
	var list []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Error("访问路径 %q 出错：%v\n", path, err)
			return err
		}
		if strings.HasSuffix(info.Name(), ext) {
			list = append(list, path)
		}
		return nil
	})

	return list, err
}

func GetFileName(filePath string) string {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return ""
	}
	return filepath.Base(filePath)
}
