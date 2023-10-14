package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fatih/color"
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
		run()
	},
}

var setEnvCommand = `export GOSS_OSS_ACCESS_KEY_ID=xxxxxxxxxxxxxxx  # your oss access_key_id
export GOSS_OSS_ACCESS_KEY_SECRET=xxxxxxxxxxxxxxxxxxx  # you oss access_key_secret
export GOSS_OSS_BUCKET_NAME=xxxxxxxx  # you oss bucket name
export GOSS_OSS_FOLDER_NAME=xxxxxx  # the folder name where you save files on OSS, example: images/avatar
export GOSS_OSS_ENDPOINT=xxxxxxxxxxxxxxxx  # you oss bucket endpoint, example: oss-cn-hongkong.aliyuncs.com`

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

	rootCmd.PersistentFlags().StringVar(&config.Dirname, "dir", ".", "The directory where the HTML file is located, defaults to the current directory")
	rootCmd.PersistentFlags().StringSliceVar(&config.Filenames, "files", nil, "your-filename. The target file to be processed. When the value is *, all HTML format files in the current directory are selected by default. If multiple files need to be selected, please use the \",\" separator, for example: \"index.html,home.html\".")
	rootCmd.PersistentFlags().BoolVar(&config.Backup, "backup", true, "Back up the original files to prevent their loss")

	rootCmd.PersistentFlags().StringVar(&config.Proxy, "proxy", "", "network proxy address")

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

func run() {
	err := initializeData()
	if err != nil {
		return
	}
	img.Run()
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

func initializeData() error {
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
	// if !utils.IsLocalPath(config.Dirname) {
	// 	utils.Warning("the value [%s] of dir is not a local path, defaults to using the current path", config.Dirname)
	// 	config.Dirname = "."
	// }
	// if !utils.IsLocalPath(config.GOSS_OSS_FOLDER_NAME) {
	// 	utils.Warning("the value [%s] of oss_folder_name is not a local path, defaults to using the root directory", config.GOSS_OSS_FOLDER_NAME)
	// 	config.GOSS_OSS_FOLDER_NAME = ""
	// }

	if len(config.Filenames) == 0 {
		filePaths, err := utils.GetDirFiles(config.Dirname, "html")
		if err != nil {
			utils.Error("Failed to obtain the list of HTML files in the %s directory", config.Dirname)
			return err
		}
		filesName, err := multipleChoice(filePaths)
		for _, fileName := range filesName {
			config.Filenames = append(config.Filenames, path.Join(config.Dirname, fileName))
		}
	} else {
		filenames := []string{}
		for _, fileName := range config.Filenames {
			var filePath string
			if utils.IsAbsolutePath(fileName) {
				filePath = fileName
			} else {
				//相对路径
				filePath = path.Join(config.Dirname, fileName)
			}
			// 检查是否存在该文件
			exist, err := utils.PathExists(filePath)
			if err != nil {
				utils.Error("Unexpected error reading %s file: %s", filePath, err)
				return err
			}
			if !exist {
				//不存在
				utils.Error("file [%s] not exist, please check the value of --files entered", filePath)
				return fmt.Errorf("file [%s] not exist, please check the value of --files entered", filePath)
			}
			filenames = append(filenames, filePath)
		}
	}
	return nil
}

func multipleChoice(options []string) ([]string, error) {
	// 初始化选择状态
	selected := make([]bool, len(options))

	// 设置颜色
	title := color.New(color.FgCyan, color.Bold)
	keyword := color.New(color.FgYellow)
	selectTxt := color.New(color.FgHiMagenta)

	if len(options) == 0 {
		keyword.Println("该目录下，没有html文件")
		return nil, nil
	}

	// 打印列表
	title.Println("请选择你需要处理的文件（用逗号分隔）:")
	keyword.Println(" 0. 全选")
	keyword.Println("-1. 退出不选")
	for i, option := range options {
		selectTxt.Printf("%d. %s\n", i+1, option)
	}

	for {
		// 获取用户选择
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入选项的编号: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			utils.Error("读取输入时发生错误:%s", err)
			return nil, err
		}

		// 解析用户选择
		choices := strings.Split(input[:len(input)-1], ",")
		valid := true
		for _, choice := range choices {
			num, err := strconv.Atoi(strings.TrimSpace(choice))
			if err != nil || num < -1 || num > len(options) {
				utils.Warning("无效的选项编号，请重新输入")
				valid = false
				break
			}
			if num == 0 {
				// 全选
				for i := range selected {
					selected[i] = true
				}
				break
			} else if num == -1 {
				// 退出不选
				return nil, nil
			} else {
				selected[num-1] = true
			}
			selected[num-1] = true
		}

		if valid {
			break
		}
	}

	choicesOption := []string{}
	// 输出用户选择
	fmt.Println("您选择了以下选项:")
	for i, isSelected := range selected {
		if isSelected {
			color.Green("- %s", options[i])
			choicesOption = append(choicesOption, options[i])
		}
	}

	return choicesOption, nil
}
