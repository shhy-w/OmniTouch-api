package router

import (
	"bytes"
	"context"
	"net/http"

	rate_limiter "git.uozi.com/uozi/rate-limiter-go"
	"git.uozi.org/uozi/cosy-example-api/internal/limiter"
	"git.uozi.org/uozi/cosy-example-api/internal/user"
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
)

// AuthRequired 函数用于实现token鉴权和更新临期token和通过token得到当前用户
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 实现token鉴权和更新临期token和通过token得到当前用户
		u, err := user.CurrentUser(c)
		if err != nil {
			// 如果鉴权失败，返回401状态码和错误信息
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
				"error":   err.Error(),
			})
			return
		}
		// 发现被禁用的用户
		if u.Status == model.UserStatusBan {
			// 如果用户被禁用，返回403状态码和错误信息
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "this user is banned",
			})
			return
		}
		// if u.use
		c.Set("user", u)
		// // 在c中新增键值对burn_id
		// c.Set("burn_id", u.burnID)
		// 获取用户权限
		c.Set("permissions_map", u.GetPermissionsMap())

		c.Next()
	}
}

func AuthAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 实现token鉴权和更新临期token和通过token得到当前用户
		u, err := user.CurrentUser(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
				"error":   err.Error(),
			})
			return
		}
		// 根据用户组只放行管理员
		if u.UserGroupID != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "forbidden",
			})
			return
		}

		// if u.use
		c.Set("user", u)
		// // 在c中新增键值对burn_id
		// c.Set("burn_id", u.burnID)
		// 获取用户权限
		c.Set("permissions_map", u.GetPermissionsMap())

		c.Next()
	}
}
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet("user").(*model.User)
		if u.UserGroupID == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "forbidden",
			})
			return
		}
		c.Next()
	}
}

func buildLimiterKey(c *gin.Context) string {
	return "limiter:" + c.FullPath() + ":" + c.ClientIP() + ":" + c.GetHeader("X-Fingerprint")
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func LimiterMiddleware(conf *rate_limiter.LimitConf) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := buildLimiterKey(c)

		lm := limiter.GetLimiter()
		result, err := lm.Allow(context.Background(), key, conf)
		if err != nil {
			return
		}
		if result.Allowed == 0 && result.Remaining == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Your operation is too frequent",
			})
			return
		}

		c.Next()

		//// 如果对所有请求都进行限流，则直接放行
		//if onlySuccess {
		//	c.Next()
		//	return
		//}

		//responseBodyWriter := &responseWriter{
		//	body:           bytes.NewBufferString(""),
		//	ResponseWriter: c.Writer,
		//}
		//c.Writer = responseBodyWriter
		//
		//c.Next()

		//respStatusCode := c.Writer.Status()
		//go func() {
		//	if respStatusCode != http.StatusNoContent && respStatusCode != http.StatusOK {
		//	}
		//}()
	}
}
