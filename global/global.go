package global

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/licat233/goss/common"
)

var (
	Bucket *oss.Bucket
)

func Initialize() error {
	var err error
	if Bucket == nil {
		Bucket, err = common.NewOssBucket()
		if err != nil {
			return err
		}
	}
	return nil
}
