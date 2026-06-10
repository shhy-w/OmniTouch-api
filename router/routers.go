package router

import (
	"net/http"

	miniappAPI "git.uozi.org/uozi/cosy-example-api/api/miniapp"
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
	miniappAPI.InitRouter(r)

}
