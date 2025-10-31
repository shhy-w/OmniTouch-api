package upload

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"git.uozi.org/uozi/cosy-example-api/internal/minio"
	"git.uozi.org/uozi/cosy-example-api/model"
	"git.uozi.org/uozi/cosy-example-api/query"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const notTrustedIdx = -1

func Put(c *gin.Context) (upload *model.Upload, err error) {
	file, err := c.FormFile("file")
	if err != nil {
		return
	}
	to := c.PostForm("to")
	idx := lo.IndexOf(GetAllowedTo(), To(to))
	if idx == notTrustedIdx {
		return nil, errors.New("target category is not trusted")
	}
	keepName := cast.ToBool(c.PostForm("keep_name"))
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return
	}
	defer os.Remove(tempDir)

	localPath := path.Join(tempDir, file.Filename)
	ext := strings.ToLower(filepath.Ext(file.Filename))

	uuidStr := uuid.New().String()
	var url string
	if keepName {
		url = path.Join(uuidStr, file.Filename)
	} else {
		url = path.Join(uuidStr + ext)
	}

	err = c.SaveUploadedFile(file, localPath)
	if err != nil {
		return
	}

	mtype, err := mimetype.DetectFile(localPath)
	var contentType string
	if err != nil {
		contentType = "application/octet-stream"
	} else {
		contentType = mtype.String()
	}

	err = minio.Put(url, localPath, contentType)
	if err != nil {
		return
	}

	var thumbnail string
	if contentType == "application/pdf" {
		thumbnail = CreatePDFThumbnail(localPath, url)
	}

	u := query.Upload
	upload = &model.Upload{
		Name:      file.Filename,
		MIME:      contentType,
		Thumbnail: thumbnail,
		Path:      url,
		Size:      cast.ToInt64(file.Size),
		To:        to,
	}

	err = u.Create(upload)
	if err != nil {
		return
	}

	err = os.Remove(localPath)
	if err != nil {
		return
	}

	return
}
