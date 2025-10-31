package global

// import (
// 	"git.uozi.org/uozi/hotel-api/api"
// 	"git.uozi.org/uozi/hotel-api/internal/user"
// 	"git.uozi.org/uozi/hotel-api/model"
// 	"git.uozi.org/uozi/hotel-api/query"
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/uozi-tech/cosy"
// 	"github.com/uozi-tech/cosy/logger"
// 	"github.com/uozi-tech/cosy/sandbox"
// 	"golang.org/x/crypto/bcrypt"
// 	"net/http"
// 	"testing"
// )

// func Test_buildTokenKey(t *testing.T) {
// 	assert.Equal(t, "token:123", api.BuildTokenKey("123"))
// 	assert.Equal(t, "token:abc", api.BuildTokenKey("abc"))
// }

// func TestLogin(t *testing.T) {
// 	sandbox.NewInstance("../../app.testing.ini", "mysql").
// 		RegisterModels(model.User{}).
// 		Run(func(instance *sandbox.Instance) {
// 			db := cosy.UseDB()
// 			query.Init(db)
// 			model.Use(db)

// 			pwd, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

// 			db.Create(&[]model.User{{
// 				Email:       "test",
// 				Name:        "test",
// 				Password:    string(pwd),
// 				Status:      model.UserStatusActive,
// 				UserGroupID: 1,
// 			}, {
// 				Email:       "banned",
// 				Name:        "banned",
// 				Password:    string(pwd),
// 				Status:      model.UserStatusBan,
// 				UserGroupID: 1,
// 			}})

// 			r := cosy.GetEngine()
// 			r.POST("/login", Login)
// 			r.DELETE("/logout", func(c *gin.Context) {
// 				_, err := user.CurrentUser(c)
// 				if err != nil {
// 					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 						"message": "unauthorized",
// 						"error":   err.Error(),
// 					})
// 					return
// 				}
// 			}, Logout)
// 			c := instance.GetClient()
// 			t.Run("success", func(t *testing.T) {
// 				resp, err := c.Post("/login", map[string]interface{}{
// 					"email":    "test",
// 					"password": "test",
// 				})
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				assert.Equal(t, http.StatusOK, resp.StatusCode)
// 				var data LoginResponse
// 				err = resp.To(&data)
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}

// 				t.Run("test logout", func(t *testing.T) {
// 					c.AddHeader("Token", data.Token)
// 					resp, err := c.Delete("/logout", nil)
// 					if err != nil {
// 						t.Error(err)
// 						return
// 					}
// 					assert.Equal(t, http.StatusNoContent, resp.StatusCode)
// 				})
// 			})

// 			t.Run("password incorrect", func(t *testing.T) {
// 				resp, err := c.Post("/login", map[string]interface{}{
// 					"email":    "test",
// 					"password": "123456",
// 				})
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
// 				var data LoginResponse
// 				err = resp.To(&data)
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				assert.Equal(t, ErrPasswordIncorrect, data.Code)

// 				resp, err = c.Post("/login", map[string]interface{}{
// 					"email":    "123456",
// 					"password": "123456",
// 				})
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
// 				err = resp.To(&data)
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				assert.Equal(t, ErrPasswordIncorrect, data.Code)
// 			})

// 			t.Run("user banned", func(t *testing.T) {
// 				resp, err := c.Post("/login", map[string]interface{}{
// 					"email":    "banned",
// 					"password": "test",
// 				})
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
// 				var data LoginResponse
// 				err = resp.To(&data)
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				assert.Equal(t, ErrUserBanned, data.Code)
// 			})

// 			t.Run("max attempts", func(t *testing.T) {
// 				for i := 0; i < 10; i++ {
// 					_, err := c.Post("/login", map[string]interface{}{
// 						"email":    "test",
// 						"password": "123456",
// 					})
// 					if err != nil {
// 						t.Error(err)
// 						return
// 					}
// 				}
// 				//time.Sleep(10 * time.Millisecond)
// 				resp, err := c.Post("/login", map[string]interface{}{
// 					"email":    "test",
// 					"password": "123456",
// 				})
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
// 				var data LoginResponse
// 				err = resp.To(&data)
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				logger.Info(data)
// 				assert.Equal(t, ErrMaxAttempts, data.Code)
// 			})
// 		})
// }
