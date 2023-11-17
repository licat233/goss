package _html

import (
	"errors"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/licat233/goss/config"
	"github.com/licat233/goss/pkg/bucket"
	"github.com/licat233/goss/utils"
)

const Name = "tools for processing HTML files"

var CheckoutFileExt = []string{"html", "htm"}

var SupportTags = []string{"img", "link", "script"}

type Html struct {
	Bucket    *bucket.Bucket
	filenames []string
	backup    bool
	tags      []string
}

func New() *Html {
	return &Html{
		Bucket:    bucket.New(config.Filenames),
		filenames: config.Filenames,
		backup:    config.Backup,
		tags:      config.HtmlTags,
	}
}

func Run() error {
	return New().Run()
}

func (s *Html) init() error {
	if !s.Bucket.Status {
		return errors.New("bucket not create")
	}
	return nil
}

func (s *Html) Run() error {
	var err error
	if err = s.init(); err != nil {
		return err
	}
	for _, filename := range s.filenames {
		if err = s.handlerSingleFile(filename); err != nil {
			return err
		}
	}
	err = s.Bucket.CreateLog(nil)
	if err != nil {
		return err
	}
	return nil
}

func (o *Html) handlerSingleFile(htmlFilePath string) error {
	//先判断该文件是否存在，root.go中的初始化阶段，已经进行检查过了，这里无需再检查
	// exist, err := utils.PathExists(htmlFilePath)
	// if err != nil {
	// 	utils.Error("Unexpected error processing %s file: %s", htmlFilePath, err)
	// 	return err
	// }
	// if !exist {
	// 	utils.Warning("The %s file does not exist", htmlFilePath)
	// 	return nil
	// }

	// 读取 HTML 文件
	htmlFile, err := os.Open(htmlFilePath)
	if err != nil {
		utils.Error("Error reading file: %s", err)
		return err
	}

	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		utils.Error("Error reading HTML file: %s", err)
		return err
	}

	hasModify := false
	for _, tag := range o.tags {
		switch tag {
		case "img":
			// 处理img标签
			isModify := o.handlerImgTag(doc, htmlFilePath)
			if !hasModify {
				hasModify = isModify
			}
		case "link":
			// 处理link标签
			isModify := o.handlerLinkTag(doc, htmlFilePath)
			if !hasModify {
				hasModify = isModify
			}
		case "script":
			// 处理script标签
			isModify := o.handlerScriptTag(doc, htmlFilePath)
			if !hasModify {
				hasModify = isModify
			}
		}
	}

	if hasModify {
		//备份
		if o.backup {
			if err := utils.BackupFile(htmlFilePath); err != nil {
				utils.Error("backup file [%s] error: %s", htmlFilePath, err)
				return err
			}
		}
		updatedHTML, err := doc.Html()
		if err != nil {
			utils.Error("unexpected error: %s", err)
			return err
		}
		err = os.WriteFile(htmlFilePath, []byte(updatedHTML), 0666)
		if err != nil {
			utils.Error("unexpected error: %s", err)
		}
	}
	// utils.Success("The image of [%s] file has been processed", htmlFilePath)
	return nil
}

// func (h *Html) existTags(searchTags []string) bool {
// 	for _, s1 := range SupportTags {
// 		for _, s2 := range searchTags {
// 			if s2 == "*" || strings.EqualFold(s1, s2) {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }
