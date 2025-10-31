package minio

// import (
// 	"os"
// 	"testing"
// 	"time"

// 	"git.uozi.org/uozi/hotel-api/settings"
// 	"github.com/google/uuid"
// 	"github.com/uozi-tech/cosy"
// 	cSettings "github.com/uozi-tech/cosy/settings"
// )

// func TestMinio(t *testing.T) {
// 	// prepare an empty temp file with uuid name
// 	name := "test-" + uuid.New().String()

// 	cSettings.Register("minio", settings.MinioSettings)

// 	go cosy.Boot("../../app.testing.ini")

// 	time.Sleep(500 * time.Millisecond)

// 	file, err := os.CreateTemp("", name)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer os.Remove(file.Name())
// 	contentType := "application/octet-stream"
// 	// Upload a file to minio
// 	err = Put(name, file.Name(), contentType)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	// Get a file from minio
// 	err = Get(name, file.Name())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	url, err := PresignedGetObject(name)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(url)
// 	// Remove a file from minio
// 	err = Remove(name)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
