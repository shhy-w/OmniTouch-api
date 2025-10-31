package settings

type OSS struct {
	AccessKeyId     string
	AccessKeySecret string
	EndPoint        string
	BucketName      string
	BaseUrl         string
}

var OssSettings = &OSS{}
