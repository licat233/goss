package bucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly"
	"github.com/licat233/goss/config"
	"github.com/licat233/goss/global"
	"github.com/licat233/goss/utils"
)

type Bucket struct {
	*oss.Bucket
	UploadedFiles map[string]string // 用于记录已上传的文件
	filenames     []string
	bucketName    string
	bucketDomain  string
	folderName    string
	endpoint      string
	Status        bool
	HasUpload     bool
}

func New(filenames []string) *Bucket {
	uploadedFiles, err := ReadLog()
	if err != nil {
		utils.Warning("read log error: %s", err)
	}
	if len(filenames) == 0 {
		filenames = config.Filenames
	}
	return &Bucket{
		Bucket:        global.Bucket,
		UploadedFiles: uploadedFiles,
		filenames:     filenames,
		bucketName:    config.GOSS_OSS_BUCKET_NAME,
		bucketDomain:  fmt.Sprintf("//%s.%s", config.GOSS_OSS_BUCKET_NAME, config.GOSS_OSS_ENDPOINT),
		folderName:    config.GOSS_OSS_FOLDER_NAME,
		endpoint:      config.GOSS_OSS_ENDPOINT,
		Status:        global.Bucket != nil,
		HasUpload:     false,
	}
}

func (s *Bucket) UploadFiles(filenames []string) error {
	if len(filenames) == 0 {
		filenames = config.Filenames
	}
	bar := pb.StartNew(len(s.filenames))
	for _, filename := range s.filenames {
		// 上传文件
		_, err := s.UploadToOss(filename)
		bar.Increment()
		if err != nil {
			utils.Warning("upload file faild: %s", err)
			return err
		}
	}
	bar.Finish()
	return nil
}

func (s *Bucket) UploadToOss(fileSrc string) (string, error) {
	// 检查路径是否为网络URL
	isNetworkURL := utils.IsURL(fileSrc)
	savePath := s.NewSaveFilePath(fileSrc)
	//检查是否已经上传过
	if fileUrl, ok := s.UploadedFiles[fileSrc]; ok {
		return fileUrl, nil
	}
	if !s.HasUpload {
		s.HasUpload = true
	}
	// 如果是本地路径，进行本地上传
	if !isNetworkURL {
		//如果是以/开头的路径，则根目录是config.Dirname，而不是
		fileSrc = path.Join(config.Dirname, fileSrc)
		err := s.PutObjectFromFile(savePath, fileSrc)
		if err != nil {
			return "", fmt.Errorf("Failed to upload local file to OSS: %s \n %s", err, fileSrc)
		}
	} else {
		//如果是网络文件
		fileUrl := fileSrc
		// 如果是当前bucket上的，则不用上传
		if s.IsOnCurrentBucket(fileSrc) {
			//不允许上传，返回原路径
			return fileSrc, nil
		}

		//先下载网络流，再上传到bucket
		body, status := s.RequestFileBody(fileUrl)
		if status != 200 {
			return "", fmt.Errorf("\nrequst 301 error: %s\n", fileUrl)
		}
		if len(body) == 0 {
			return "", fmt.Errorf("\nrequst file failed: %s\n", fileUrl)
		}

		if err := s.PutObject(savePath, bytes.NewReader(body)); err != nil {
			return "", err
		}
	}
	newSrc, err := url.JoinPath(s.bucketDomain, savePath)
	if err != nil {
		return "", fmt.Errorf("unexpected error: %s", err)
	}
	defer func() {
		err = s.CreateLog(nil)
		if err != nil {
			utils.Error("create log error: %s", err)
		}
	}()
	s.UploadedFiles[fileSrc] = newSrc
	return newSrc, nil
}

func (s *Bucket) RequestFileBody(fileUrl string) (body []byte, status int) {
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

func (s *Bucket) NewSaveFilePath(fileSrc string) string {
	saveFilename := utils.UUIDhex()
	ext := utils.FileExt(fileSrc)
	if ext == "" {
		ext = "other"
	}
	ext = strings.TrimLeft(ext, ".")
	saveFilename = strings.TrimRight(saveFilename, ".")
	saveFilename = saveFilename + "." + ext
	savePath := path.Join(s.folderName, saveFilename)
	return savePath
}

func (s *Bucket) IsOnCurrentBucket(fileUrl string) bool {
	u, err := url.Parse(fileUrl)
	if err != nil {
		//说明不是url，肯定不是
		return false
	}
	if u.Host == "" {
		//不是url，肯定不是
		return false
	}
	host := fmt.Sprintf("%s.%s", s.bucketName, s.endpoint)
	contain := strings.Contains(u.Host, host)
	return contain
}

type fileLog struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

func (s *Bucket) CreateLog(uploadedFiles map[string]string) error {
	if len(uploadedFiles) == 0 {
		uploadedFiles = s.UploadedFiles
	}
	fileLogs := []fileLog{}
	for k, v := range uploadedFiles {
		fileLogs = append(fileLogs, fileLog{
			Path: k,
			Url:  v,
		})
	}
	// 将数据转换为 JSON 字符串
	jsonStr, err := json.MarshalIndent(fileLogs, "", "  ")
	if err != nil {
		return err
	}
	// 将 JSON 字符串写入文件
	err = os.WriteFile("goss_upload_log.json", jsonStr, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadLog() (uploadedFiles map[string]string, err error) {
	uploadedFiles = make(map[string]string)
	file, err := os.Open("goss_upload_log.json")
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var fileLogs []fileLog
	err = decoder.Decode(&fileLogs)
	if err != nil {
		return
	}

	for _, fileLog := range fileLogs {
		uploadedFiles[fileLog.Path] = fileLog.Url
	}
	return
}
