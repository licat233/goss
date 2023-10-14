package utils

import "path/filepath"

// func IsLocalPath(pathName string) bool {
// 	if pathName == "." {
// 		return true
// 	}
// 	isAbsolutePath := filepath.IsAbs(pathName)
// 	baseName := filepath.Base(pathName)
// 	isRelativePath := (baseName != pathName)
// 	return isAbsolutePath || isRelativePath
// }

// func IsRelativePath(pathName string) bool {
// 	if pathName == "." {
// 		return true
// 	}
// 	baseName := filepath.Base(pathName)
// 	return (baseName != pathName)
// }

func IsAbsolutePath(pathName string) bool {
	return filepath.IsAbs(pathName)
}
