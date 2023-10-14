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
		defer func() {
			utils.Success("done.")
		}()
		img.Run()
	},
}

var setEnvCommand = `export GOSS_OSS_ACCESS_KEY_ID=xxxxxxxxxxxxxxx  # your oss access_key_id
export GOSS_OSS_ACCESS_KEY_SECRET=xxxxxxxxxxxxxxxxxxx  # you oss access_key_secret
export GOSS_OSS_BUCKET_NAME=xxxxxxxx  # you oss bucket name
export GOSS_OSS_FOLDER_NAME=xxxxxx  # the folder name where you save files on OSS, example: images/avatar
export GOSS_OSS_ENDPOINT=xxxxxxxxxxxxxxxx  # you oss bucket endpoint, example: oss-cn-hongkong.aliyuncs.com
`

var setEnvCommandWithColor = fmt.Sprintf("\033[33m%s\033[0m\033[32m", setEnvCommand)

var rootCmd = &cobra.Command{
	Use:        "goss",
	Aliases:    []string{},
	SuggestFor: []string{},
	Short:      "This is a tool for uploading images from files to OSS",
	GroupID:    "",
	Long:       fmt.Sprintf("Upload the images in the specified HTML file to OSS and update the relevant image links.\ncurrent version: %s\nGithub: https://github.com/licat233/goss.\nif you want to set nev: \n%s", config.ProjectVersion, setEnvCommandWithColor),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s requires at least one argument", cmd.CommandPath())
		}
		return nil
	},
	Version: config.ProjectVersion,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var IsDev bool

func init() {
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_ENDPOINT, "endpoint", "", "your-oss-endpoint. Default use of environment variable value of GOSS_OSS_ENDPOINT, example: oss-cn-hongkong.aliyuncs.com")
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_ACCESS_KEY_ID, "id", "", "your-access-key-id. Default use of environment variable value of GOSS_OSS_ACCESS_KEY_ID")
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_ACCESS_KEY_SECRET, "secret", "", "your-access-key-secret. Default use of environment variable value of GOSS_OSS_ACCESS_KEY_SECRET")
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_BUCKET_NAME, "bucket", "", "your-bucket-name. Default use of environment variable value of GOSS_OSS_BUCKET_NAME")
	rootCmd.PersistentFlags().StringVar(&config.GOSS_OSS_FOLDER_NAME, "folder", "", "your-oss-folder. Default use of environment variable value of GOSS_OSS_FOLDER_NAME")

	rootCmd.PersistentFlags().StringSliceVar(&config.Filenames, "files", []string{"index.html"}, "your-filename. The target file to be processed. When the value is *, all HTML format files in the current directory are selected by default. If multiple files need to be selected, please use the \",\" separator, for example: \"index.html,home.html\".")

	rootCmd.PersistentFlags().BoolVar(&IsDev, "dev", false, "dev mode, print error message")

	rootCmd.SetHelpTemplate(greenColorizeHelp(rootCmd.HelpTemplate()))

	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(startCmd)
}

func greenColorizeHelp(template string) string {
	// 在这里使用ANSI转义码添加颜色
	// 参考ANSI转义码文档：https://en.wikipedia.org/wiki/ANSI_escape_code
	// 将标题文本设置为绿色
	template = "\033[34m" + template + "\033[0m"
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
	initializeData()
	if err := rootCmd.Execute(); err != nil {
		utils.Warning(err.Error())
		os.Exit(1)
	}
}

func initializeData() {
	if config.GOSS_OSS_ENDPOINT == "" {
		config.GOSS_OSS_ENDPOINT = os.Getenv("GOSS_OSS_ENDPOINT")
	}
	if config.GOSS_OSS_ACCESS_KEY_ID == "" {
		config.GOSS_OSS_ACCESS_KEY_ID = os.Getenv("GOSS_OSS_ACCESS_KEY_ID")
	}
	if config.GOSS_OSS_ACCESS_KEY_SECRET == "" {
		config.GOSS_OSS_ACCESS_KEY_SECRET = os.Getenv("GOSS_OSS_ACCESS_KEY_SECRET")
	}
	if config.GOSS_OSS_BUCKET_NAME == "" {
		config.GOSS_OSS_BUCKET_NAME = os.Getenv("GOSS_OSS_BUCKET_NAME")
	}
	if config.GOSS_OSS_FOLDER_NAME == "" {
		config.GOSS_OSS_FOLDER_NAME = os.Getenv("GOSS_OSS_FOLDER_NAME")
	}
}
