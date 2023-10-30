package _html

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gocolly/colly"
	"github.com/licat233/goss/config"
	"github.com/licat233/goss/global"
	"github.com/licat233/goss/utils"
)

const Name = "tools for processing HTML files"

var CheckoutFileExt = []string{"html", "htm"}

var SupportTags = []string{"img", "link", "script"}

type Html struct {
	bucket         *oss.Bucket
	uploadedImages map[string]string // 用于记录已上传的文件
	filenames      []string
	bucketName     string
	bucketDomain   string
	folderName     string
	endpoint       string
	Status         bool
	backup         bool
	tags           []string
}

func New() *Html {
	return &Html{
		bucket:         global.Bucket,
		uploadedImages: map[string]string{},
		filenames:      config.Filenames,
		bucketName:     config.GOSS_OSS_BUCKET_NAME,
		bucketDomain:   fmt.Sprintf("//%s.%s", config.GOSS_OSS_BUCKET_NAME, config.GOSS_OSS_ENDPOINT),
		folderName:     config.GOSS_OSS_FOLDER_NAME,
		endpoint:       config.GOSS_OSS_ENDPOINT,
		Status:         false,
		backup:         config.Backup,
		tags:           config.HtmlTags,
	}
}

func Run() error {
	return New().Run()
}

func (s *Html) init() error {
	if s.bucket == nil {
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
	return nil
}

func (o *Html) isOnCurrentBucket(imgUrl string) bool {
	u, err := url.Parse(imgUrl)
	if err != nil {
		//说明不是url，肯定不是
		return false
	}
	if u.Host == "" {
		//不是url，肯定不是
		return false
	}
	host := fmt.Sprintf("%s.%s", o.bucketName, o.endpoint)
	contain := strings.Contains(u.Host, host)
	return contain
}

func (o *Html) newSaveFilePath(fileSrc string) string {
	saveFilename := utils.UUIDhex()
	ext := utils.FileExt(fileSrc)
	if ext == "" {
		ext = "other"
	}
	ext = strings.TrimLeft(ext, ".")
	saveFilename = strings.TrimRight(saveFilename, ".")
	saveFilename = saveFilename + "." + ext
	savePath := path.Join(o.folderName, saveFilename)
	return savePath
}

func (o *Html) requestFileBody(fileUrl string) (body []byte, status int) {
	// 发送HTTP GET请求获取文件内容
	c := colly.NewCollector()
	c.SetRequestTimeout(10 * time.Second)
	if proxyURL := strings.TrimSpace(config.Proxy); proxyURL != "" {
		if err := c.SetProxy(proxyURL); err != nil {
			utils.Error("set proxy error: %s", err)
			return
		}
	}

	c.OnResponse(func(r *colly.Response) {
		status = r.StatusCode
		if r.StatusCode != 200 {
			return
		}
		body = r.Body
	})
	c.Visit(fileUrl)
	c.Wait()
	return
}

func (o *Html) uploadToOss(fileSrc string) (string, error) {
	// 检查路径是否为网络URL
	isNetworkURL := utils.IsURL(fileSrc)
	savePath := o.newSaveFilePath(fileSrc)

	//检查是否已经上传过
	if imgUrl, ok := o.uploadedImages[fileSrc]; ok {
		return imgUrl, nil
	}

	// 如果是本地路径，进行本地上传
	if !isNetworkURL {
		//如果是以/开头的路径，则根目录是config.Dirname，而不是
		fileSrc = path.Join(config.Dirname, fileSrc)
		err := o.bucket.PutObjectFromFile(savePath, fileSrc)
		if err != nil {
			return "", fmt.Errorf("Failed to upload local file to OSS: %s \n %s", err, fileSrc)
		}
	} else {
		//如果是网络文件
		imageURL := fileSrc
		// 如果是当前bucket上的，则不用上传
		if o.isOnCurrentBucket(imageURL) {
			//不允许上传，返回原路径
			return fileSrc, nil
		}

		//先下载网络流，再上传到bucket
		body, status := o.requestFileBody(imageURL)
		if status != 200 {
			return "", fmt.Errorf("\nrequst 301 error: %s\n", imageURL)
		}
		if len(body) == 0 {
			return "", fmt.Errorf("\nrequst file failed: %s\n", imageURL)
		}

		if err := o.bucket.PutObject(savePath, bytes.NewReader(body)); err != nil {
			return "", err
		}
	}
	newSrc, err := url.JoinPath(o.bucketDomain, savePath)
	if err != nil {
		return "", fmt.Errorf("unexpected error: %s", err)
	}
	o.uploadedImages[fileSrc] = newSrc
	return newSrc, nil
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
