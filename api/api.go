package api

import (
	"net/http"
	"net/url"
	"time"

	"git.uozi.org/uozi/cosy-example-api/internal/acl"
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/redis"
)

func CurrentUser(c *gin.Context) *model.User {
	return c.MustGet("user").(*model.User)
}

func Can(c *gin.Context, subject acl.Subject, action acl.Action) bool {
	permissionsMap := c.MustGet("permissions_map").(acl.Map)
	return acl.Can(permissionsMap, subject, action)
}

func Guard(subject acl.Subject, actions ...acl.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissionsMap := c.MustGet("permissions_map").(acl.Map)
		// logger.Debug(permissionsMap, actions)
		// Iterate through the list of actions, if any action is allowed, go next.
		for _, action := range actions {
			if acl.Can(permissionsMap, subject, action) {
				c.Next()
				return
			}
		}

		// If none of the actions are allowed, return 403 Forbidden.
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "resource limited by acl rules",
		})
	}
}

func CosyGuard[T any](c cosy.ICurd[T], subject acl.Subject) {
	c.BeforeCreate(Guard(subject, acl.Write)).
		BeforeModify(Guard(subject, acl.Write)).
		BeforeGet(Guard(subject, acl.Read)).
		BeforeGetList(Guard(subject, acl.Read)).
		BeforeDestroy(Guard(subject, acl.Write)).
		BeforeRecover(Guard(subject, acl.Write))
}

func BuildTokenKey(token string) string {
	return "token:" + token
}

func GenerateToken(user *model.User) (string, error) {
	token, err := acl.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	// 设置token在redis中的过期时间
	expiration := time.Hour * 6
	// UserGroupID != 0说明是管理员,令牌有效期延长
	if user.UserGroupID != 0 {
		expiration = time.Hour * 24 * 15
	}
	// 用token生成唯一的key,用户id作为value.
	err = redis.Set(BuildTokenKey(token), user.ID, expiration)
	return token, err
}

func FileAttachmentFromBytes(c *gin.Context, data []byte, filename string) {
	c.Writer.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	c.Data(http.StatusOK, "application/octet-stream", data)
}
