package router

import (
	"net/http"

	burn_history "git.uozi.org/uozi/cosy-example-api/api/admin/burned_history"
	"git.uozi.org/uozi/cosy-example-api/api/global"
	"git.uozi.org/uozi/cosy-example-api/api/user"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

func InitRouter() {
	r := cosy.GetEngine()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok321",
		})
	})
	global.InitRouter(r)
	// adminRouter := r.Group("admin", AuthAdminRequired(), AdminRequired())
	adminRouter := r.Group("admin", AuthAdminRequired())
	{
		adminRouter.GET("/user", global.GetCurrentUser)
		adminRouter.POST("/user", global.UpdateCurrentUser)
		user.InitUserRouter(adminRouter)
		user.InitUserGroupRouter(adminRouter)
		burn_history.InitBurnHistoryRouter(adminRouter)
	}

}
