package upload

import (
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"git.uozi.org/uozi/cosy-example-api/internal/minio"
	"github.com/gen2brain/go-fitz"
	"github.com/uozi-tech/cosy/logger"
)

func CreatePDFThumbnail(file string, objectName string) (url string) {
	url = filepath.Base(strings.TrimSuffix(objectName, filepath.Ext(file)) + ".jpg")
	doc, err := fitz.New(file)
	if err != nil {
		logger.Error(err)
		return
	}
	defer doc.Close()

	img, err := doc.ImageDPI(0, 96.0)
	if err != nil {
		logger.Error(err)
		return
	}

	f, err := os.CreateTemp("", url)
	if err != nil {
		logger.Error(err)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()

	err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		logger.Error(err)
		return
	}

	err = minio.Put(url, f.Name(), "image/jpeg")
	if err != nil {
		logger.Error(err)
		return
	}

	return
}
