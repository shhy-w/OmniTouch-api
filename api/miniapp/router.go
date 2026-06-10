package miniapp

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/auth/wx-login", WxLogin)

	devices := api.Group("/devices", AuthRequired())
	{
		devices.GET("", ListDevices)
		devices.POST("/bind", BindDevice)
		devices.GET("/:device_id", GetDevice)
		devices.PATCH("/:device_id", RenameDevice)
		devices.DELETE("/:device_id", DeleteDevice)
		devices.POST("/:device_id/commands", SendCommand)
		devices.POST("/:device_id/voice", SendVoice)
		devices.GET("/:device_id/status", GetStatus)
		devices.GET("/:device_id/events", ListEvents)
	}
}
