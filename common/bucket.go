package common

import (
	"errors"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/licat233/goss/config"
	"github.com/licat233/goss/utils"
)

func NewOssBucket() (*oss.Bucket, error) {
	// 初始化 OSS 客户端
	endpoint := config.GOSS_OSS_ENDPOINT
	if endpoint == "" {
		utils.Warning("Please configure the value of GOSS_OSS_ENDPOINT, use the '%s - h' command to view help information", config.ProjectName)
		return nil, errors.New("Missing core configuration")
	}
	accessKeyID := config.GOSS_OSS_ACCESS_KEY_ID
	if accessKeyID == "" {
		utils.Warning("Please configure the value of GOSS_OSS_ACCESS_KEY_ID, use the '%s - h' command to view help information", config.ProjectName)
		return nil, errors.New("Missing core configuration")
	}
	accessKeySecret := config.GOSS_OSS_ACCESS_KEY_SECRET
	if accessKeySecret == "" {
		utils.Warning("Please configure the value of GOSS_OSS_ACCESS_KEY_SECRET, use the '%s - h' command to view help information", config.ProjectName)
		return nil, errors.New("Missing core configuration")
	}
	bucketName := config.GOSS_OSS_BUCKET_NAME
	if bucketName == "" {
		utils.Warning("Please configure the value of GOSS_OSS_BUCKET_NAME, use the '%s - h' command to view help information", config.ProjectName)
		return nil, errors.New("Missing core configuration")
	}

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		utils.Error("Error creating OSS client: %s", err)
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		utils.Error("Error accessing OSS bucket: %s", err)
		return nil, err
	}
	return bucket, nil
}
