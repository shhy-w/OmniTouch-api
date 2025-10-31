package oss

// import (
// 	"bytes"
// 	"encoding/json"
// 	"git.uozi.org/uozi/hotel-api/model"
// 	"git.uozi.org/uozi/hotel-api/query"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/spf13/cast"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/uozi-tech/cosy"
// 	"github.com/uozi-tech/cosy/router"
// 	"github.com/uozi-tech/cosy/sandbox"
// 	"io"
// 	"mime/multipart"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"path"
// 	"strings"
// 	"testing"
// )

// func TestOSS(t *testing.T) {
// 	sandbox.NewInstance("../../app.testing.ini", "mysql").
// 		RegisterModels(model.User{}, model.Upload{}).
// 		Run(func(instance *sandbox.Instance) {

// 			db := cosy.UseDB()
// 			query.Init(db)

// 			cosy.GetEngine().POST("/upload", func(c *gin.Context) {
// 				dir := path.Join("test", uuid.NewString())

// 				o, err := NewOSS()
// 				if err != nil {
// 					cosy.ErrHandler(c, err)
// 					return
// 				}

// 				if c.Query("is_client") == "true" {
// 					c.Set("user", &model.User{
// 						Model: model.Model{
// 							ID: 1,
// 						},
// 						UserGroupID: 0,
// 					})
// 				} else {
// 					c.Set("user", &model.User{
// 						Model: model.Model{
// 							ID: 1,
// 						},
// 						UserGroupID: 1,
// 					})
// 				}

// 				res, err := o.UploadSingleFile(c, dir)
// 				if err != nil {
// 					cosy.ErrHandler(c, err)
// 					return
// 				}

// 				// new record to the database
// 				upload := &model.Upload{
// 					UserId: 1,
// 					Path:   res.Url,
// 					Size:   res.Size,
// 					MIME:   res.MIME,
// 					Name:   res.Filename,
// 					To:     "",
// 				}

// 				u := query.Upload
// 				if err = u.Create(upload); err != nil {
// 					cosy.ErrHandler(c, err)
// 					return
// 				}

// 				ok, err := o.IsObjectExist(res.Url)
// 				if err != nil {
// 					cosy.ErrHandler(c, err)
// 					return
// 				}
// 				assert.Equal(t, true, ok)

// 				err = o.DeleteObject(res.Url)
// 				if err != nil {
// 					cosy.ErrHandler(c, err)
// 					return
// 				}

// 				c.JSON(http.StatusOK, upload)
// 			})

// 			fileName := "test-" + uuid.NewString() + ".txt"
// 			file, err := os.CreateTemp("", fileName)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			defer os.Remove(file.Name())

// 			testUpload(t, fileName, file, true, false)
// 			testUpload(t, fileName, file, false, false)
// 			testUpload(t, fileName, file, true, true)
// 		})

// }

// func testUpload(t *testing.T, name string, file *os.File, keepName, isClient bool) {
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	part, err := writer.CreateFormFile("file", name)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	_, err = io.Copy(part, file)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	_ = writer.WriteField("keep_name", cast.ToString(keepName))

// 	err = writer.Close()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	req := httptest.NewRequest("POST", "/upload?is_client="+cast.ToString(isClient), body)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())

// 	w := httptest.NewRecorder()
// 	router.GetEngine().ServeHTTP(w, req)
// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	var respData struct {
// 		Path string `json:"path"`
// 	}

// 	respBytes, _ := io.ReadAll(resp.Body)

// 	_ = json.Unmarshal(respBytes, &respData)

// 	urlS := strings.Split(respData.Path, "/")

// 	if keepName {
// 		assert.Equal(t, name, urlS[len(urlS)-1])
// 	} else {
// 		assert.NotEqual(t, name, urlS[len(urlS)-1])
// 	}
// }
