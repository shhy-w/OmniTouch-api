package global

import (
	"net/http"

	"git.uozi.org/uozi/cosy-example-api/api"
	"git.uozi.org/uozi/cosy-example-api/internal/user"
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/logger"
)

// 删
func DestroyUser(c *gin.Context) {
	cosy.Core[model.User](c).Destroy()
}

func DestroyUser_(c *gin.Context) {

	cosy.Core[model.User](c).PermanentlyDelete()
}

// TODO:此处有 hotel
func ModifyUser(c *gin.Context) {
	logger.Debug(c.Request)
	core := cosy.Core[model.User](c).SetValidRules(gin.H{
		"name":     "omitempty",
		"email":    "omitempty",
		"password": "omitempty",
		"hotel_id": "omitempty",
		"status":   "omitempty",
		// ... 其他字段
	})

	core.BeforeExecuteHook(user.EncryptPassword).
		SetNextHandler(GetUser).Modify()
}

func GetUser(c *gin.Context) {
	cosy.Core[model.User](c).Get()
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "ok",
	// })
}

func GetAllUser(c *gin.Context) {
	cosy.Core[model.User](c).PagingList()
	// logger.Debug("返回所有用户")
	// data, err := cosy.Core[model.User](c).ListAllData()
	// if !err  {
	// 	logger.Error("查询所有用户错误")
	// 	// cosy.ErrHandler(c)
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": data,
	// })
}

func GetCurrentUser(c *gin.Context) {
	user := api.CurrentUser(c)
	c.JSON(http.StatusOK, user)
}

func UpdateCurrentUser(c *gin.Context) {
	var info struct {
		AvatarID interface{} `json:"avatar_id"`
		Name     string      `json:"name"`
		Phone    string      `json:"phone"`
		Email    string      `json:"email"`
	}
	if !cosy.BindAndValid(c, &info) {
		return
	}

	user := api.CurrentUser(c)
	db := cosy.UseDB(c)

	err := db.Model(user).Updates(&model.User{
		Name:  info.Name,
		Email: info.Email,
	}).Error
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
