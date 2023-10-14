package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	err := filepath.WalkDir(dir, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			Error("访问路径 %q 出错：%v\n", path, err)
			return err
		}
		if dirEntry.IsDir() && path != dir {
			// 如果是目录且不是当前目录，则跳过
			return filepath.SkipDir
		}
		if strings.HasSuffix(dirEntry.Name(), ext) {
			list = append(list, path)
		}
		return nil
	})

	return list, err
}

// 获取文件名
func GetFileName(filePath string) string {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return ""
	}
	return filepath.Base(filePath)
}

// 批量获取文件名
func GetFilesName(filesPath ...string) (filesName []string) {
	for _, filePath := range filesPath {
		filesName = append(filesName, GetFileName(filePath))
	}
	return
}

// 备份文件
func BackupFile(srcPath string) error {
	// 获取源文件所在目录和文件名
	srcDir := filepath.Dir(srcPath)
	srcFileName := filepath.Base(srcPath)

	// 构建备份文件名
	timestamp := time.Now().Unix()
	fileNameWithoutExt := strings.TrimSuffix(srcFileName, filepath.Ext(srcFileName))
	destFileName := fmt.Sprintf("%s-%d%s", fileNameWithoutExt, timestamp, filepath.Ext(srcFileName))
	destPath := filepath.Join(srcDir, destFileName)

	// 读取源文件
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// 写入备份文件
	err = os.WriteFile(destPath, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
