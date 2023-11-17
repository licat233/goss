package utils

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 获取指定目录下的所有文件或目录
// 如果exts为空，则获取所有文件和目录
// 如果exts包含"*"，则获取所有文件
// 如果exts不包含"*"，但包含其他字符，则获取对应格式名的文件
func GetAllFiles(dirPath string, exts []string) ([]string, error) {
	if dirPath == "" {
		dirPath = "."
	}
	isAll := len(exts) == 0
	justContainCurrentDir := false
	if SliceContain(exts, "*") {
		isAll = true
		justContainCurrentDir = true
		exts = nil
	}

	//处理文件格式名
	for i := range exts {
		ext := exts[i]
		if ext == "" {
			continue
		}
		ext = strings.TrimLeft(ext, ".")
		if ext != "" {
			ext = "." + ext
		}
		exts[i] = ext
	}

	var list []string
	err := filepath.WalkDir(dirPath, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			Error("访问路径 %q 出错：%v\n", path, err)
			return err
		}

		if dirEntry.IsDir() {
			if path == dirPath {
				return nil
			}
			//如果只包含当前目录，则忽略其它目录
			//前提是不能是当前目录
			if justContainCurrentDir {
				return filepath.SkipDir
			}
			return nil
		}
		if isAll {
			//如果全部选择则不用判断格式
			list = append(list, path)
			return nil
		}
		if SliceContain(exts, filepath.Ext(dirEntry.Name())) {
			//全部选择
			list = append(list, path)
		}
		return nil
	})

	return list, err
}

// 获取指定目录下的某些格式的文件列表
func GetDirFiles(dir string, exts []string) ([]string, error) {
	if dir == "" {
		dir = "."
	}
	isAll := len(exts) == 0
	if !isAll {
		isAll = SliceContain(exts, "*")
		exts = SliceRemoves(exts, []string{"*"})
	}

	if !isAll {
		//处理文件格式名
		for i := range exts {
			ext := exts[i]
			if ext == "" {
				continue
			}
			ext = strings.TrimLeft(ext, ".")
			if ext != "" {
				ext = "." + ext
			}
			exts[i] = ext
		}
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
		if isAll || SliceContain(exts, filepath.Ext(dirEntry.Name())) {
			//全部选择
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

func getFilesAndDirs(dirPath string) ([]string, error) {
	var files []string

	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		files = append(files, fileInfo.Name())
	}

	return files, nil
}
