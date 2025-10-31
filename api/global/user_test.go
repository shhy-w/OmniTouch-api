package global

// import (
// 	"git.uozi.org/uozi/hotel-api/model"
// 	"git.uozi.org/uozi/hotel-api/query"
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/uozi-tech/cosy"
// 	"github.com/uozi-tech/cosy/sandbox"
// 	"golang.org/x/crypto/bcrypt"
// 	"testing"
// )

// func TestUser(t *testing.T) {
// 	sandbox.NewInstance("../../app.testing.ini", "mysql").
// 		RegisterModels(model.User{}, model.Upload{}).Run(func(instance *sandbox.Instance) {
// 		r := cosy.GetEngine()
// 		pwd, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
// 		testUser := &model.User{
// 			Email:       "test",
// 			Name:        "test",
// 			UserGroupID: 0,
// 			Password:    string(pwd),
// 			Status:      model.UserStatusActive,
// 		}

// 		r.GET("/user", func(c *gin.Context) {
// 			c.Set("user", testUser)
// 		}, GetCurrentUser)
// 		r.POST("/user", func(c *gin.Context) {
// 			c.Set("user", testUser)
// 		}, UpdateCurrentUser)

// 		db := cosy.UseDB()
// 		query.Init(db)
// 		db.Create(testUser)
// 		db.Create(&[]model.Upload{
// 			{
// 				Model: model.Model{
// 					ID: 1,
// 				},
// 				UserId: 1,
// 				Path:   "abc.png",
// 			},
// 			{
// 				Model: model.Model{
// 					ID: 2,
// 				},
// 				UserId: 1,
// 				Path:   "bcd.png",
// 			},
// 		})

// 		t.Run("test get current user", func(t *testing.T) {
// 			c := instance.GetClient()
// 			resp, err := c.Get("/user")
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			var user model.User
// 			err = resp.To(&user)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			assert.Equal(t, testUser.Email, user.Email)
// 			assert.Equal(t, testUser.Name, user.Name)
// 			assert.Equal(t, testUser.Phone, user.Phone)
// 			assert.Equal(t, testUser.AvatarID, user.AvatarID)
// 			assert.Contains(t, user.Avatar.Path, "abc.png")
// 			assert.Equal(t, testUser.Status, user.Status)
// 			assert.Equal(t, testUser.UserGroupID, user.UserGroupID)
// 		})

// 		t.Run("test update current user", func(t *testing.T) {
// 			c := instance.GetClient()
// 			resp, err := c.Post("/user", map[string]interface{}{
// 				"avatar_id":     2,
// 				"name":          "test1",
// 				"phone":         "test1",
// 				"email":         "test1",
// 				"user_group_id": 1,
// 			})
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			var user model.User
// 			err = resp.To(&user)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			assert.Equal(t, "test1", user.Email)
// 			assert.Equal(t, "test1", user.Name)
// 			assert.Equal(t, "test1", user.Phone)
// 			assert.Equal(t, testUser.AvatarID, user.AvatarID)
// 			assert.Contains(t, user.Avatar.Path, "bcd.png")
// 			assert.Equal(t, testUser.Status, user.Status)
// 			assert.Equal(t, testUser.UserGroupID, user.UserGroupID)
// 		})
// 	})
// }
