package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/licat233/goss/config"
	"github.com/licat233/goss/utils"
)

var colorId int = 30

func getCmdColor() int {
	colorId += 1
	if colorId < 30 || (colorId > 37 && colorId < 90) || colorId > 97 {
		return 30
	}
	return colorId
}

func setColorizeHelp(template string) string {
	// 在这里使用ANSI转义码添加颜色
	// 参考ANSI转义码文档：https://en.wikipedia.org/wiki/ANSI_escape_code
	return fmt.Sprintf("\033[%dm%s\033[0m", getCmdColor(), template)
}

func setColorPart(part string) string {
	return fmt.Sprintf("\033[%dm%s", getCmdColor(), part)
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

	if err := checkHtmlConfig(); err != nil {
		return err
	}

	if utils.SliceContain(config.Filenames, "*") {
		config.Filenames = nil
	}

	return nil
}

func checkoutFiles(fileExts []string) ([]string, error) {
	//如果不指定文件，或者指定文件为*，则以文件格式作为依据
	if len(config.Filenames) == 0 || utils.SliceContain(config.Filenames, "*") {
		//如果用户设置了文件格式，则以用户的设置为准
		if len(config.Exts) != 0 {
			fileExts = config.Exts
		}
		filePaths, err := utils.GetAllFiles(config.Dirname, fileExts)
		if err != nil {
			utils.Error("Failed to obtain the list of HTML files in the %s directory", config.Dirname)
			return nil, err
		}
		filesName, err := multipleSelectionFiles(filePaths)
		for _, fileName := range filesName {
			config.Filenames = append(config.Filenames, path.Join(config.Dirname, fileName))
		}
	} else {
		//检查指定的文件
		for i, fileName := range config.Filenames {
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
				return nil, err
			}
			if !exist {
				//不存在
				utils.Error("file [%s] not exist, please check the value of --files entered", filePath)
				return nil, fmt.Errorf("file [%s] not exist, please check the value of --files entered", filePath)
			}
			config.Filenames[i] = filePath
		}
	}
	config.Filenames = utils.SliceRemoves(config.Filenames, []string{"*", ""})
	return config.Filenames, nil
}

func multipleSelectionFiles(options []string) ([]string, error) {
	// 初始化选择状态
	selected := make([]bool, len(options))

	// 设置颜色
	title := color.New(color.FgCyan, color.Bold)
	keyword := color.New(color.FgYellow)
	selectTxt := color.New(color.FgHiMagenta)

	if len(options) == 0 {
		dirname, err := filepath.Abs(config.Dirname)
		if err != nil {
			utils.Error("无法获取到你设置的目录: dir = %s \n error: %s", dirname, err.Error())
			return nil, err
		}
		keyword.Printf("目录下没有文件，请确认目标目录: %s\n", dirname)
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
		title.Print("请输入选项的编号: ")
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
	title.Println("您选择了以下文件:")
	for i, isSelected := range selected {
		if isSelected {
			color.Green("- %s", options[i])
			choicesOption = append(choicesOption, options[i])
		}
	}

	return choicesOption, nil
}
