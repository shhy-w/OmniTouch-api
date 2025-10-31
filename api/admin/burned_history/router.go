package burn_history

import (
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

// 可供参考
// 初始化销毁路由
func InitBurnHistoryRouter(r *gin.RouterGroup) {

	r.GET("/burned_history", GetList)
	r.GET("/burned_history/:id", GetBurnedHistory)
	r.POST("/burned_history/:id", UpdateBurnedHistory)
}

func GetList(c *gin.Context) {
	core := cosy.Core[model.BurnedHistory](c).
		SetFussy("client_ip", "en_key").
		SetIn("status").
		SetPreloads("FirmwareFile")
	core.PagingList()
}

func UpdateBurnedHistory(c *gin.Context) {
	core := cosy.Core[model.BurnedHistory](c).SetValidRules(gin.H{
		"remark": "required",
	})

	core.SetNextHandler(GetBurnedHistory).Modify()
}
func GetBurnedHistory(c *gin.Context) {
	core := cosy.Core[model.BurnedHistory](c)
	core.Get()
}
