package user

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	// "fmt"
	"time"

	"git.uozi.org/uozi/cosy-example-api/internal/acl"
	"git.uozi.org/uozi/cosy-example-api/internal/helper"
	"git.uozi.org/uozi/cosy-example-api/model"
	"git.uozi.org/uozi/cosy-example-api/query"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/redis"
)

// GetByID get user by id (with redis cache)
func GetByID(id uint64) (user *model.User, err error) {
	key := helper.BuildUserKey(id)
	userStr, err := redis.Get(key)
	if err != nil || userStr == "" {
		u := query.User
		user, err = u.Preload(u.UserGroup).FirstByID(id)

		if err != nil {
			return
		}
		// hide password
		user.Password = ""
		bytes, _ := json.Marshal(user)
		_ = redis.Set(key, string(bytes), 5*time.Minute)

		user.UserGroup, err = user.GetUserGroup()

		return
	}

	user = &model.User{}
	err = json.Unmarshal([]byte(userStr), user)

	if err != nil {
		return nil, err
	}

	user.UserGroup, err = user.GetUserGroup()

	return
}

func CurrentToken(c *gin.Context) (token string) {
	// 首先检查请求头中是否包含 "Token" 字段。
	if len(c.Request.Header["Token"]) == 0 {
		// 如果请求头中没有 "Token"，则尝试从查询参数中获取 "token"。
		if c.Query("token") == "" {
			return ""
		}
		// 对查询参数的"token"进行解码。
		tmp, _ := base64.StdEncoding.DecodeString(c.Query("token"))
		token = string(tmp)
	} else {
		// 若请求头有 "Token" 字段，则直接取第一个值作为token。
		token = c.Request.Header["Token"][0]
	}
	// 返回当前请求体的token
	return token
}

// CurrentUser get current user from token
func CurrentUser(c *gin.Context) (u *model.User, err error) {
	// validate token
	token := CurrentToken(c)

	if token == "" {
		return nil, errors.New("token not found")
	}

	var claims *acl.JWTClaims
	claims, err = acl.ValidateJWT(token)
	if err != nil {
		return
	}

	// get user by id
	u, err = GetByID(claims.UserID)
	if err != nil {
		return
	}

	_, err = redis.Get("token:" + token)
	if err != nil {
		logger.Error("token not found in redis", token)
		return nil, err
	}
	expiresAt := redis.TTL("token:" + token).Seconds()

	if expiresAt <= 0 {
		logger.Error("token expired", token)
		return nil, errors.New("token expired")
	}

	if expiresAt < 30*60 {
		err = redis.Set("token:"+token, u.ID, 1*time.Minute)
		if err != nil {
			return nil, err
		}
		// refresh token
		newToken, err := acl.GenerateJWT(u.ID)
		if err != nil {
			return nil, err
		}
		err = redis.Set("token:"+newToken, u.ID, 3*time.Hour)
		if err != nil {
			return nil, err
		}
		c.Header("refresh-token", newToken)
	}

	u.UpdateLastActive()

	logger.Info("[Current User]", u.Name)

	return
}
