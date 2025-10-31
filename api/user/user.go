package user

import (
	"net/http"

	"git.uozi.org/uozi/cosy-example-api/internal/user"
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

func InitUserRouter(g *gin.RouterGroup) {
	c := cosy.Api[model.User]("users")
	// api.CosyGuard[model.User](c, acl.User)
	// 新增-密码加密
	c.CreateHook(func(c *cosy.Ctx[model.User]) {
		c.BeforeDecodeHook(user.EncryptPassword)
	})
	// 修改-密码加密
	c.ModifyHook(func(c *cosy.Ctx[model.User]) {
		c.BeforeDecodeHook(user.EncryptPassword)
	})

	// 查询列表-活跃时间
	c.BeforeGetList(func(c *gin.Context) {
		user.PersistLastActive(c)
	})

	// 查询详情-活跃时间
	//删除-超级管理员不能删除
	c.DestroyHook(func(c *cosy.Ctx[model.User]) {

		c.BeforeExecuteHook(func(ctx *cosy.Ctx[model.User]) {
			if ctx.OriginModel.ID == 1 {
				ctx.JSON(http.StatusNotAcceptable, gin.H{
					"message": "Cannot delete the super admin",
				})

				ctx.Abort()
			}
		})

	})

	c.InitRouter(g)
}
