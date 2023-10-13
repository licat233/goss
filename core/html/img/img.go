package img

import (
	"bytes"
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
	"github.com/licat233/goss/utils"
)

type Img struct {
	bucket         *oss.Bucket
	uploadedImages map[string]string // 用于记录已上传的图片
	filenames      []string
	bucketName     string
	bucketDomain   string
	folderName     string
	endpoint       string
	Status         bool
}

func New() *Img {
	i := &Img{
		bucket:         nil,
		uploadedImages: map[string]string{},
		filenames:      config.Filenames,
		bucketName:     config.GOSS_OSS_BUCKET_NAME,
		bucketDomain:   fmt.Sprintf("//%s.%s", config.GOSS_OSS_BUCKET_NAME, config.GOSS_OSS_ENDPOINT),
		folderName:     config.GOSS_OSS_FOLDER_NAME,
		endpoint:       config.GOSS_OSS_ENDPOINT,
		Status:         false,
	}
	i.Status = i.init() == nil
	return i
}

func Run() {
	New().Run()
}

func (s *Img) init() error {
	bucket, err := newOssBucket()
	if err != nil {
		return err
	}
	s.bucket = bucket
	return nil
}

func (s *Img) Run() {
	if !s.Status {
		return
	}
	var err error
	filenames := s.filenames
	if len(s.filenames) > 0 && s.filenames[0] == "*" {
		//遍历当前目录下的所有html文件
		filenames, err = utils.GetDirFiles(".", ".html")
		if err != nil {
			return
		}
	}
	for _, filename := range filenames {
		if err = s.handlerSingleFile(filename); err != nil {
			return
		}
	}
}

func (o *Img) allowUpload(imgUrl string) bool {
	//判断是否是已经上传到当前bucket的图片
	u, err := url.Parse(imgUrl)
	if err != nil {
		//说明不是url，可以上传
		return true
	}
	if u.Host == "" {
		//不是url，可上传
		return true
	}
	host := fmt.Sprintf("%s.%s", o.bucketName, o.endpoint)
	contain := strings.Contains(u.Host, host)
	if contain {
		//已经存在了，不用上传
		return false
	}
	return true
}

func (o *Img) newSaveFilePath(imagSrc string) string {
	saveFilename := utils.UUID()
	if ext := utils.FileExt(imagSrc); ext != "" {
		saveFilename = saveFilename + ext
	}
	savePath := path.Join(o.folderName, saveFilename)
	return savePath
}

func (o *Img) getImageBody(imageURL string) (body []byte, status int) {
	// 发送HTTP GET请求获取图片内容
	c := colly.NewCollector()
	c.SetRequestTimeout(10 * time.Second)
	// proxyURL := "http://127.0.0.1:7890"
	// if err := c.SetProxy(proxyURL); err != nil {
	// 	log.Fatalln(err)
	// }
	c.OnResponse(func(r *colly.Response) {
		status = r.StatusCode
		if r.StatusCode != 200 {
			return
		}
		body = r.Body
	})
	c.Visit(imageURL)
	c.Wait()
	return
}

func (o *Img) uploadToOss(imagSrc string) (string, error) {
	// 检查图片路径是否为网络URL
	isNetworkURL := utils.IsURL(imagSrc)
	savePath := o.newSaveFilePath(imagSrc)

	// 如果是本地图片，进行本地上传
	if !isNetworkURL {
		if imgUrl, ok := o.uploadedImages[imagSrc]; ok {
			return imgUrl, nil
		}
		err := o.bucket.PutObjectFromFile(savePath, imagSrc)
		if err != nil {
			return "", fmt.Errorf("Failed to upload local image to OSS: %s \n %s", err, imagSrc)
		}
	} else {
		imageURL := imagSrc
		// 如果是网络图片，进行网络流文件上传
		if !o.allowUpload(imageURL) {
			//不允许上传，返回原路径
			return imagSrc, nil
		}

		body, status := o.getImageBody(imageURL)
		if status != 200 {
			return "", fmt.Errorf("\nrequst 301 error: %s\n", imageURL)
		}
		if len(body) == 0 {
			return "", fmt.Errorf("\nrequst image failed: %s\n", imageURL)
		}

		if err := o.bucket.PutObject(savePath, bytes.NewReader(body)); err != nil {
			return "", err
		}
	}
	newSrc, err := url.JoinPath(o.bucketDomain, savePath)
	if err != nil {
		return "", fmt.Errorf("unexpected error: %s", err)
	}
	o.uploadedImages[imagSrc] = newSrc
	return newSrc, nil
}

func (o *Img) handlerSingleFile(htmlFilePath string) error {
	//先判断该文件是否存在
	exist, err := utils.PathExists(htmlFilePath)
	if err != nil {
		utils.Error("Unexpected error processing %s file: %s", htmlFilePath, err)
		return err
	}
	if !exist {
		utils.Warning("The %s file does not exist", htmlFilePath)
		return nil
	}

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

	// 遍历并处理img标签
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}
		src = strings.TrimSpace(src)
		if src == "" {
			return
		}
		newSrc, err := o.uploadToOss(src)
		if err != nil {
			utils.Warning("upload image faild: %s", err)
			return
		}
		// 更新img标签的src属性
		s.SetAttr("src", newSrc)

		src, exists = s.Attr("data-src")
		if !exists {
			return
		}
		newSrc, err = o.uploadToOss(src)
		if err != nil {
			utils.Warning("upload image faild: %s", err)
			return
		}
		// 更新img标签的src属性
		s.SetAttr("data-src", newSrc)
	})

	updatedHTML, err := doc.Html()
	if err != nil {
		utils.Error("unexpected error: %s", err)
		return err
	}

	err = os.WriteFile(htmlFilePath, []byte(updatedHTML), 0666)
	if err != nil {
		utils.Error("unexpected error: %s", err)
	}

	utils.Success("The image of [%s] file has been processed", htmlFilePath)
	return nil
}
