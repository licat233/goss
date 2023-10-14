# goss

将指定的html文件内的图片上传到oss中，并更新链接

## install

```shell
go install github.com/licat233/goss@latest
```

## help

```shell
goss -h
```

## upgrade

```shell
goss upgrade
```

## configure OSS

```shell
export GOSS_OSS_ACCESS_KEY_ID=xxxxxxxxxxxxxxx  # your oss access_key_id
export GOSS_OSS_ACCESS_KEY_SECRET=xxxxxxxxxxxxxxxxxxx  # you oss access_key_secret
export GOSS_OSS_BUCKET_NAME=xxxxxxxx  # you oss bucket name
export GOSS_OSS_FOLDER_NAME=xxxxxx  # the folder name where you save files on OSS, example: images/avatar
export GOSS_OSS_ENDPOINT=xxxxxxxxxxxxxxxx  # you oss bucket endpoint, example: oss-cn-hongkong.aliyuncs.com
```
