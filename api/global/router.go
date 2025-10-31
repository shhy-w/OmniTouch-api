package global

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	r.POST("admin/login", Login)
	r.DELETE("admin/logout", Logout)

}
