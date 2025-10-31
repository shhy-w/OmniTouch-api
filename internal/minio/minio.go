package minio

import (
	"context"
	"strings"
	"time"

	"git.uozi.org/uozi/cosy-example-api/settings"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/uozi-tech/cosy/redis"
)

var serverBaseUrl = ""

func initClient() (*minio.Client, error) {
	endpoint := settings.MinioSettings.Endpoint
	accessKeyID := settings.MinioSettings.AccessKeyID
	secretAccessKey := settings.MinioSettings.AccessKeySecret

	if settings.MinioSettings.Secure {
		serverBaseUrl = "https://"
	} else {
		serverBaseUrl = "http://"
	}

	serverBaseUrl += endpoint

	// Initialize minio client object.
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: settings.MinioSettings.Secure,
	})
}

func Put(objectName, filePath, contentType string) (err error) {
	ctx := context.Background()
	bucketName := settings.MinioSettings.BucketName
	minioClient, err := initClient()
	if err != nil {
		return err
	}
	// Upload the test file with FPutObject
	_, err = minioClient.FPutObject(ctx, bucketName, objectName, filePath,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	return
}

func Get(objectName, filePath string) (err error) {
	ctx := context.Background()
	bucketName := settings.MinioSettings.BucketName
	minioClient, err := initClient()
	if err != nil {
		return err
	}
	// Download the object to a file
	return minioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
}

func Remove(objectName string) (err error) {
	ctx := context.Background()
	minioClient, err := initClient()

	if err != nil {
		return
	}

	err = minioClient.RemoveObject(ctx, settings.MinioSettings.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return
	}

	return
}

func PresignedGetObject(objectName string) (url string, err error) {
	url, err = getCacheOrPresignedURL(objectName)
	return
}

func getCacheOrPresignedURL(objectName string) (url string, err error) {
	var sb strings.Builder
	sb.WriteString("minio:")
	sb.WriteString(objectName)
	url, err = redis.Get(sb.String())
	if err != nil || url == "" {
		ctx := context.Background()
		minioClient, err := initClient()
		if err != nil {
			return "", err
		}

		presignedURL, err := minioClient.PresignedGetObject(ctx, settings.MinioSettings.BucketName,
			objectName, 24*time.Hour, nil)
		if err != nil {
			return "", err
		}

		// url = strings.ReplaceAll(presignedURL.String(), serverBaseUrl, "")
		url := presignedURL.String()
		_ = redis.Set(sb.String(), url, 23*time.Hour)
		return url, nil
	}
	return
}
