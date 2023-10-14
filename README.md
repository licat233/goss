# goss

将指定的html文件内的图片上传到oss中，并更新链接

## install

```shell
go install github.com/licat233/goss@latest
```

## help

```shell
$goss -h
Upload the images in the specified HTML file to OSS and update the relevant image links.
current version: v1.0.6
Github: https://github.com/licat233/goss.
if you want to set nev:
export GOSS_OSS_ACCESS_KEY_ID=xxxxxxxxxxxxxxx  # your oss access_key_id
export GOSS_OSS_ACCESS_KEY_SECRET=xxxxxxxxxxxxxxxxxxx  # you oss access_key_secret
export GOSS_OSS_BUCKET_NAME=xxxxxxxx  # you oss bucket name
export GOSS_OSS_FOLDER_NAME=xxxxxx  # the folder name where you save files on OSS, example: images/avatar
export GOSS_OSS_ENDPOINT=xxxxxxxxxxxxxxxx  # you oss bucket endpoint, example: oss-cn-hongkong.aliyuncs.com

Usage:
  goss [flags]
  goss [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       run goss
  upgrade     Upgrade goss to latest version

Flags:
      --backup            Back up the original files to prevent their loss (default true)
      --bucket string     your-bucket-name. Default use of environment variable value of GOSS_OSS_BUCKET_NAME
      --dev               dev mode, print error message
      --dir string        The directory where the HTML file is located, defaults to the current directory (default ".")
      --endpoint string   your-oss-endpoint. Default use of environment variable value of GOSS_OSS_ENDPOINT, example: oss-cn-hongkong.aliyuncs.com
      --files strings     your-filename. The target file to be processed. When the value is *, all HTML format files in the current directory are selected by default. If multiple files need to be selected, please use the "," separator, for example: "index.html,home.html".
      --folder string     your-oss-folder. Default use of environment variable value of GOSS_OSS_FOLDER_NAME
  -h, --help              help for goss
      --id string         your-access-key-id. Default use of environment variable value of GOSS_OSS_ACCESS_KEY_ID
      --proxy string      network proxy address
      --secret string     your-access-key-secret. Default use of environment variable value of GOSS_OSS_ACCESS_KEY_SECRET
  -v, --version           version for goss

Use "goss [command] --help" for more information about a command.
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
