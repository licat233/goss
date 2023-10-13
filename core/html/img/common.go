package img

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/licat233/goss/config"
	"github.com/licat233/goss/utils"
)

func newOssBucket() (*oss.Bucket, error) {
	// 初始化 OSS 客户端
	endpoint := config.GOSS_OSS_ENDPOINT
	accessKeyID := config.GOSS_OSS_ACCESS_KEY_ID
	accessKeySecret := config.GOSS_OSS_ACCESS_KEY_SECRET
	bucketName := config.GOSS_OSS_BUCKET_NAME
	// folderName := config.GOSS_OSS_FOLDER_NAME

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
