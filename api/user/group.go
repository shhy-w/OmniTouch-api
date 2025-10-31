package user

import (
	"net/http"

	"git.uozi.org/uozi/cosy-example-api/api"
	"git.uozi.org/uozi/cosy-example-api/internal/acl"
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy"
)

func removeCachedUserGroup(c *cosy.Ctx[model.UserGroup]) {
	c.Model.CleanCache()
}

func InitUserGroupRouter(r *gin.RouterGroup) {
	c := cosy.Api[model.UserGroup]("user_groups")
	api.CosyGuard[model.UserGroup](c, acl.User)

	c.DestroyHook(func(c *cosy.Ctx[model.UserGroup]) {
		c.BeforeExecuteHook(func(ctx *cosy.Ctx[model.UserGroup]) {
			if cast.ToInt(ctx.Param("id")) == 1 {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "Delete default group is forbidden.",
				})
				ctx.Abort()
			}
		})
		c.ExecutedHook(removeCachedUserGroup)
	})

	c.ModifyHook(func(c *cosy.Ctx[model.UserGroup]) {
		c.BeforeExecuteHook(func(ctx *cosy.Ctx[model.UserGroup]) {
			if ctx.OriginModel.ID == 1 {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "Modify default group is forbidden.",
				})
				ctx.Abort()
			}
		})
		c.ExecutedHook(removeCachedUserGroup)
	})

	c.DestroyHook(func(ctx *cosy.Ctx[model.UserGroup]) {
		ctx.ExecutedHook(removeCachedUserGroup)
	})

	c.InitRouter(r)

	r.GET("/user_groups/subjects", api.Guard(acl.All, acl.Read, acl.Write), GetSubjects)
}

func GetSubjects(c *gin.Context) {
	// api return
	c.JSON(http.StatusOK, acl.GetSubjects())
}
