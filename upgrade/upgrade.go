package upgrade

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/licat233/goss/utils"

	"github.com/licat233/goss/config"
)

func Upgrade() {
	currentVersion := config.ProjectVersion
	latestVersion := getLatestVersion()

	if strings.EqualFold(currentVersion, latestVersion) {
		utils.Success("当前版本[%s]已是最新.", currentVersion)
		return
	}

	utils.Success("新版本[%s]可用!", latestVersion)
	if err := updateSelf(latestVersion); err != nil {
		utils.Error("更新失败：%s\n", err)
		os.Exit(1)
	}
	utils.Success("更新成功!")
	// 重新启动程序
	os.Exit(0)
}

// 获取github项目的最新release版本号
func GetLatestReleaseVersion(projectInfoURL string) (string, error) {
	command := fmt.Sprintf("wget -qO- -t1 -T2 \"%s\" | grep \"tag_name\" | head -n 1 | awk -F \":\" '{print $2}' | sed 's/\\\"//g;s/,//g;s/ //g'", projectInfoURL)
	out, err := utils.ExecShell(command)
	out = strings.TrimSpace(out)
	return out, err
}

// 获取最新版本号
func getLatestVersion() string {
	v, err := GetLatestReleaseVersion(config.ProjectInfoURL)
	if err != nil {
		utils.Error("获取最新版本号失败：%s\n", err)
		os.Exit(1)
	}
	return v
}

// 自我升级
func updateSelf(latestVersion string) error {
	goBinary, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("未找到go命令：%w", err)
	}
	url := strings.ReplaceAll(config.ProjectURL, "http://", "")
	url = strings.ReplaceAll(url, "https://", "") + "@" + latestVersion
	utils.Success("更新命令: go install %s\n正在更新...", url)
	// 构建并安装最新版本的程序
	if err := exec.Command(goBinary, "install", url).Run(); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
