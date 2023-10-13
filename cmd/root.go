package cmd

import (
	"fmt"
	"os"

	"github.com/licat233/goss/config"
	"github.com/licat233/goss/core/html/img"
	"github.com/licat233/goss/upgrade"
	"github.com/licat233/goss/utils"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number of " + config.ProjectName,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Success("current version: " + config.ProjectVersion)
	},
}

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"up", "u"},
	Short:   "Upgrade " + config.ProjectName + " to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		upgrade.Upgrade()
	},
}

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run"},
	Short:   "run " + config.ProjectName,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var rootCmd = &cobra.Command{
	Use:        "goss",
	Aliases:    []string{},
	SuggestFor: []string{},
	Short:      "This is a tool for uploading images from files to OSS",
	GroupID:    "",
	Long:       fmt.Sprintf("Upload the images in the specified HTML file to OSS and update the relevant image links.\ncurrent version: %s\nGithub: https://github.com/licat233/goss", config.ProjectVersion),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s requires at least one argument", cmd.CommandPath())
		}
		return nil
	},
	Version: config.ProjectVersion,
	Run: func(cmd *cobra.Command, args []string) {
		img.Run()
	},
}

var IsDev bool

func init() {
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_ENDPOINT, "endpoint", "", "your-oss-endpoint. Default use of environment variable value of GOSS_OSS_ENDPOINT.")
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_ACCESS_KEY_ID, "id", "", "your-access-key-id. Default use of environment variable value of GOSS_OSS_ACCESS_KEY_ID.")
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_ACCESS_KEY_SECRET, "secret", "", "your-access-key-secret. Default use of environment variable value of GOSS_OSS_ACCESS_KEY_SECRET.")
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_BUCKET_NAME, "bucket", "", "your-bucket-name. Default use of environment variable value of GOSS_OSS_BUCKET_NAME.")
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_FOLDER_NAME, "folder", "", "your-oss-folder. Default use of environment variable value of GOSS_OSS_FOLDER_NAME.")

	rootCmd.PersistentFlags().StringSliceVar(&config.Filenames, "files", []string{"index.html"}, "your-filename. The target file to be processed. When the value is *, all HTML format files in the current directory are selected by default. If multiple files need to be selected, please use the \",\" separator, for example: \"index.html,home.html\".")

	rootCmd.PersistentFlags().BoolVar(&IsDev, "dev", false, "dev mode, print error message")

	rootCmd.SetHelpTemplate(greenColorizeHelp(rootCmd.HelpTemplate()))
}

func greenColorizeHelp(template string) string {
	// 在这里使用ANSI转义码添加颜色
	// 参考ANSI转义码文档：https://en.wikipedia.org/wiki/ANSI_escape_code
	// 将标题文本设置为绿色
	template = "\033[32m" + template + "\033[0m"
	return template
}

func Execute() {
	defer func() {
		if !IsDev {
			if err := recover(); err != nil {
				utils.Warning(fmt.Sprintf("%v", err))
			}
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		utils.Warning(err.Error())
		os.Exit(1)
	}
}
