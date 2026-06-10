package miniapp

import (
	"errors"
	"strings"
	"time"

	"git.uozi.org/uozi/cosy-example-api/internal/acl"
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/redis"
)

const tokenPrefix = "miniapp_token:"

func GenerateToken(user *model.MiniappUser) (string, error) {
	token, err := acl.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	if err = redis.Set(tokenPrefix+token, user.ID, 15*24*time.Hour); err != nil {
		return "", err
	}
	return token, nil
}

func CurrentUser(c *gin.Context) (*model.MiniappUser, error) {
	token := bearerToken(c)
	if token == "" {
		token = c.GetHeader("Token")
	}
	if token == "" {
		return nil, errors.New("token not found")
	}
	claims, err := acl.ValidateJWT(token)
	if err != nil {
		return nil, err
	}
	if _, err = redis.Get(tokenPrefix + token); err != nil {
		return nil, errors.New("token expired")
	}
	var user model.MiniappUser
	if err = cosy.UseDB(c).First(&user, claims.UserID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func bearerToken(c *gin.Context) string {
	value := c.GetHeader("Authorization")
	if value == "" {
		return ""
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(value)
}
