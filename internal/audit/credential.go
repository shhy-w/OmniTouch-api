package audit

import (
	"git.uozi.org/uozi/cosy-example-api/settings"
	sls "github.com/aliyun/aliyun-log-go-sdk"
)

func getCredentialsProvider() *sls.StaticCredentialsProvider {
	accessKeyId := settings.SLSSettings.AccessKeyId
	accessKeySecret := settings.SLSSettings.AccessKeySecret

	return sls.NewStaticCredentialsProvider(accessKeyId,
		accessKeySecret, "")
}
