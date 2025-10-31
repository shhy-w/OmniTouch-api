package settings

import (
	"git.uozi.org/uozi/cosy-example-api/api"
	"git.uozi.org/uozi/cosy-example-api/internal/acl"
	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.RouterGroup) {
	r := g.Group("", api.Guard(acl.All, acl.Write))

	r.GET("settings/:name", GetSetting)
	r.POST("settings/:name", UpdateSetting)

	r.GET("settings/auth", GetAuthSettings)
	r.POST("settings/auth", SaveAuthSettings)
	r.GET("settings/auth/banned_ips", GetBanLoginIP)
	r.DELETE("settings/auth/banned_ip", RemoveBannedIP)
}
