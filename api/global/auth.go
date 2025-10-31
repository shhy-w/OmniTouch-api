package global

import (
	"net/http"
	"time"

	"git.uozi.org/uozi/cosy-example-api/api"
	"git.uozi.org/uozi/cosy-example-api/internal/runtime_settings"
	"git.uozi.org/uozi/cosy-example-api/internal/user"
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/redis"
)

type loginJson struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code"`
	Token   string      `json:"token,omitempty"`
	User    *model.User `json:"user,omitempty"`
}

func Login(c *gin.Context) {
	// s := logger.NewSessionLogger(c)
	lock, err := redis.ObtainLock("login:"+c.RemoteIP()+c.GetHeader("X-Fingerprint"), 10*time.Millisecond, nil)

	if err != nil {
		c.JSON(http.StatusTooManyRequests, user.ErrMaxAttempts)
		return
	}
	defer lock.Release(c)

	var loginFailedKey = "login_failed:" + c.RemoteIP()

	auth := runtime_settings.GetAuthSettings()

	countStr, err := redis.Get(loginFailedKey)
	if err != nil {
		_ = redis.Set(loginFailedKey, 0,
			time.Duration(
				lo.Max([]int{auth.BanThresholdMinutes, 1}),
			)*time.Minute)
	}
	failedCount := cast.ToInt(countStr)

	if auth.MaxAttempts > 0 && failedCount >= auth.MaxAttempts {
		c.JSON(http.StatusNotAcceptable, user.ErrMaxAttempts)
		return
	}

	var login loginJson
	if !cosy.BindAndValid(c, &login) {
		_, _ = redis.Incr(loginFailedKey)
		return
	}

	u, err := user.Login(c, login.Email, login.Password)
	if err != nil {
		cosy.ErrHandler(c, err)
		_, _ = redis.Incr(loginFailedKey)
		return
	}

	logger.Info("[User Login]", u.Name)

	u.UpdateLastActive()

	token, err := api.GenerateToken(u)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}

	// api return
	c.JSON(http.StatusOK, LoginResponse{
		Message: "ok",
		Token:   token,
		User:    u,
	})
}

func Logout(c *gin.Context) {
	token := user.CurrentToken(c)
	_ = redis.Del(api.BuildTokenKey(token))
	c.JSON(http.StatusNoContent, nil)
}
