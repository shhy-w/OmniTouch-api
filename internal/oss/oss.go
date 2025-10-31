package oss

import (
	"errors"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"git.uozi.org/uozi/cosy-example-api/api"
	"git.uozi.org/uozi/cosy-example-api/settings"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	// "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	// "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/logger"
)

type OSS struct {
	Client *oss.Client
	Bucket *oss.Bucket

	ExpireAt time.Duration // 过期时间
}

type UploadResult struct {
	SignedUrl string
	OssUrl    string
	Size      int64
	MIME      string
	Filename  string
}

var (
	ErrClientUploadFileType = errors.New("only accept jpeg, png, gif and webp image")
)

func NewOSS() (*OSS, error) {
	endpoint := settings.OssSettings.EndPoint
	// 确保使用 HTTPS 协议
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	client, err := oss.New(endpoint,
		settings.OssSettings.AccessKeyId, settings.OssSettings.AccessKeySecret)

	if err != nil {
		return nil, err
	}
	var bucket *oss.Bucket

	bucket, err = client.Bucket(settings.OssSettings.BucketName)
	if err != nil {
		return nil, err
	}

	return &OSS{
		Client: client,
		Bucket: bucket,
	}, nil
}

// NewOSSWithExpire 支持设置过期时间
func NewOSSWithExpire(expireSecond uint64) (*OSS, error) {
	endpoint := settings.OssSettings.EndPoint
	// 确保使用 HTTPS 协议
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	client, err := oss.New(endpoint,
		settings.OssSettings.AccessKeyId, settings.OssSettings.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	var bucket *oss.Bucket

	bucket, err = client.Bucket(settings.OssSettings.BucketName)
	if err != nil {
		return nil, err
	}

	return &OSS{
		Client:   client,
		Bucket:   bucket,
		ExpireAt: time.Duration(expireSecond) * time.Second,
	}, nil
}

// 修改上传方法，添加过期时间选项
func (o *OSS) PutObjectWithExpire(ossPath, localPath string) (string, error) {
	// 上传文件到OSS
	logger.Debug("PutObjectWithExpire->ossPath: ", ossPath)
	if err := o.Bucket.PutObjectFromFile(ossPath, localPath); err != nil {
		return "", err
	}
	// 生成带过期时间的签名URL
	signedURL, err := o.Bucket.SignURL(ossPath, oss.HTTPGet, int64(o.ExpireAt.Seconds()))
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

// 将本地文件上传到OSS
func (o *OSS) PutObject(localPath, ossPath string) error {
	return o.Bucket.PutObjectFromFile(ossPath, localPath)
}

func (o *OSS) GetObjectToFile(ossPath, localPath string, options ...oss.Option) error {
	return o.Bucket.GetObjectToFile(ossPath, localPath, options...)
}

func (o *OSS) DeleteObject(ossPath string) error {
	return o.Bucket.DeleteObject(ossPath)
}

func (o *OSS) CopyObject(src, dst string) error {
	_, err := o.Bucket.CopyObject(src, dst)
	return err
}

func (o *OSS) IsObjectExist(ossPath string) (bool, error) {
	return o.Bucket.IsObjectExist(ossPath)
}

// GenerateSignedURL 为指定的OSS对象生成新的签名URL
func (o *OSS) GenerateSignedURL(ossPath string, expireSeconds int64) (string, error) {
	return o.Bucket.SignURL(ossPath, oss.HTTPGet, expireSeconds)
}

func (o *OSS) UploadSingleFile(c *gin.Context, ossDir string) (res UploadResult, err error) {
	var file *multipart.FileHeader
	file, err = c.FormFile("file")
	if err != nil {
		return
	}

	// keep filename?
	keepFileName := cast.ToBool(c.PostForm("keep_name"))

	res.Filename = filepath.Base(file.Filename)
	res.Size = file.Size

	user := api.CurrentUser(c)

	// 10MB limit for client users
	if user.UserGroupID == 0 && file.Size > 10*1024*1024 {
		err = errors.New("file size exceeds the limit of 10MB")
		return
	}

	// get extension
	ext := filepath.Ext(file.Filename)

	// build a temporary directory
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return
	}

	defer os.RemoveAll(tempDir)

	localPath := path.Join(tempDir, file.Filename)

	var ossPath string

	if keepFileName {
		ossPath = path.Join(ossDir, uuid.New().String(), file.Filename)
	} else {
		ossPath = path.Join(ossDir, uuid.New().String()+ext)
	}

	err = c.SaveUploadedFile(file, localPath)
	if err != nil {
		return
	}

	logger.Debug("UploadSingleFile->ossPath: ", ossPath)
	logger.Debug("UploadSingleFile->ossDir: ", ossDir)

	mtype, err := mimetype.DetectFile(localPath)
	if err != nil {
		res.MIME = "application/octet-stream"
	} else {
		res.MIME = mtype.String()
	}

	err = o.PutObject(localPath, ossPath)
	if err != nil {
		return
	}

	// signedURL, err := o.PutObjectWithExpire(ossPath, localPath)
	// if err != nil {
	// 	return
	// }

	// delete the local file
	err = os.Remove(localPath)
	if err != nil {
		return
	}

	res.OssUrl = ossPath
	// res.SignedUrl = signedURL

	return
}

func (o *OSS) SignedURL(ossPath string, expireSeconds int64) (string, error) {
	return o.Bucket.SignURL(ossPath, oss.HTTPGet, expireSeconds)
}

// func (o *OSS) TestURLAccess(url string) bool {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		logger.Errorf("❌ URL access test failed: %v", err)
// 		return false
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == 200 {
// 		logger.Infof("✅ URL is accessible: %d", resp.StatusCode)
// 		return true
// 	} else {
// 		logger.Infof("❌ URL returned status: %d", resp.StatusCode)
// 		return false
// 	}
// }
