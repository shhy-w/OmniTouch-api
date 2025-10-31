package settings

import (
	"net/http"

	"git.uozi.org/uozi/cosy-example-api/internal/runtime_settings"
	"git.uozi.org/uozi/cosy-example-api/model"
	"git.uozi.org/uozi/cosy-example-api/query"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

func GetSetting(c *gin.Context) {
	name := c.Param("name")
	initialType := c.Query("type")
	s := query.Setting

	setting, _ := s.Where(s.Name.Eq(name)).First()
	if setting == nil {

		if initialType == "array" {
			c.JSON(http.StatusOK, []struct{}{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, setting.Unmarshal())
}

func UpdateSetting(c *gin.Context) {
	name := c.Param("name")
	s := query.Setting

	var json interface{}

	if !cosy.BindAndValid(c, &json) {
		return
	}

	setting, _ := s.Where(s.Name.Eq(name)).First()

	if setting == nil {
		setting = &model.Setting{
			Name: name,
		}
		setting.Marshal(&json)
		setting.Insert()
	} else {
		setting.Marshal(&json)
		_ = setting.Save()
	}

	runtime_settings.Set(name, setting.Meta)

	GetSetting(c)
}
