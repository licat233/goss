package config

const (
	ProjectName    = "goss"
	ProjectVersion = "v1.1.0-beta.6"
	ProjectInfoURL = "https://api.github.com/repos/licat233/" + ProjectName + "/releases/latest"
	ProjectURL     = "https://github.com/licat233/" + ProjectName
)

var (
	GOSS_OSS_ENDPOINT          string //Endpoint（地域节点）: oss-cn-guangzhou.aliyuncs.com
	GOSS_OSS_ACCESS_KEY_ID     string
	GOSS_OSS_ACCESS_KEY_SECRET string
	GOSS_OSS_BUCKET_NAME       string
	// GOSS_OSS_BUCKET_DOMAIN     string = os.Getenv("GOSS_OSS_BUCKET_DOMAIN") //Bucket 域名: 	licat-storage.oss-cn-guangzhou.aliyuncs.com
	GOSS_OSS_FOLDER_NAME string
)

var (
	Proxy     string   //网络代理
	Filenames []string //选择的文件
	Exts      []string //选择的文件格式
	Dirname   string   //需要读取的目录
	Backup    bool     //备份原文件，防止原文件丢失
)

var (
	HtmlTags []string //选择处理的html标签
)

// func GetVersion() string {
// 	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
// 	if err != nil {
// 		utils.Error("获取git tags出错:%s", err)
// 		return "v1.0.0"
// 	}
// 	version := strings.TrimSpace(string(out))
// 	return version
// }
