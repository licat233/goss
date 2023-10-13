package config

import (
	"os/exec"
	"strings"

	"github.com/licat233/goss/utils"
)

var (
	ProjectName    = "goss"
	ProjectVersion = GetVersion()
	ProjectInfoURL = "https://api.github.com/repos/licat233/" + ProjectName + "/releases/latest"
	ProjectURL     = "https://github.com/licat233/" + ProjectName
)

var (
	GOSS_OSS_ENDPOINT          string //Endpoint（地域节点）: oss-cn-guangzhou.aliyuncs.com
	GOSS_OSS_ACCESS_KEY_ID     string
	GOSS_OSS_ACCESS_KEY_SECRET string
	GOSS_OSS_BUCKET_NAME       string
	// GOSS_OSS_BUCKET_DOMAIN     string = os.Getenv("GOSS_OSS_BUCKET_DOMAIN") //Bucket 域名: 	licat-storage.oss-cn-guangzhou.aliyuncs.com
	GOSS_OSS_FOLDER_NAME string
)

var (
	Filenames []string
)

func GetVersion() string {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		utils.Error("获取git tags出错:%s", err)
		return "v1.0.0"
	}
	version := strings.TrimSpace(string(out))
	return version
}
