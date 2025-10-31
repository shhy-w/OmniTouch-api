package settings

import (
	"net/http"
	"strings"
	"time"

	"git.uozi.org/uozi/cosy-example-api/internal/runtime_settings"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/redis"
)

func GetAuthSettings(c *gin.Context) {
	auth := runtime_settings.GetAuthSettings()
	c.JSON(http.StatusOK, auth)
}

func SaveAuthSettings(c *gin.Context) {
	var json runtime_settings.Auth
	if !cosy.BindAndValid(c, &json) {
		return
	}
	runtime_settings.UpdateAuthSettings(&json)

	GetAuthSettings(c)
}

type BanIP struct {
	IP        string `json:"ip"`
	Attempts  int    `json:"attempts"`
	ExpiredAt int64  `json:"expired_at"`
}

func GetBanLoginIP(c *gin.Context) {
	keys, err := redis.Keys("login_failed:*")
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}

	bannedIPs := lo.Map(keys, func(ip string, key int) *BanIP {
		return &BanIP{
			IP: ip,
		}
	})
	auth := runtime_settings.GetAuthSettings()
	bannedIPs = lo.Filter(bannedIPs, func(banIP *BanIP, index int) bool {
		ttl := redis.TTL(banIP.IP)
		attemptStr, _ := redis.Get(banIP.IP)
		attampts := cast.ToInt(attemptStr)

		if ttl > 0 && attampts >= auth.MaxAttempts {
			banIP.ExpiredAt = time.Now().Add(ttl).Unix()
			banIP.Attempts = attampts
			banIP.IP = strings.TrimPrefix(banIP.IP, "login_failed:")
			return true
		}

		return false
	})

	c.JSON(http.StatusOK, bannedIPs)
}

func RemoveBannedIP(c *gin.Context) {
	var json struct {
		IP string `json:"ip"`
	}
	if !cosy.BindAndValid(c, &json) {
		return
	}
	err := redis.Del("login_failed:" + json.IP)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
