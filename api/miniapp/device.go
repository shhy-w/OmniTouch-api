package miniapp

import (
	"net/http"
	"strconv"

	internalMiniapp "git.uozi.org/uozi/cosy-example-api/internal/miniapp"
	"git.uozi.org/uozi/cosy-example-api/internal/omnitouch"
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

type bindDeviceRequest struct {
	DeviceID string `json:"device_id"`
	Name     string `json:"name"`
	SSID     string `json:"ssid"`
	Password string `json:"password"`
}

type renameDeviceRequest struct {
	Name string `json:"name"`
}

type commandRequest struct {
	Cmd    string         `json:"cmd"`
	Params map[string]any `json:"params"`
}

type voiceRequest struct {
	Text string `json:"text"`
}

func ListDevices(c *gin.Context) {
	user := mustMiniappUser(c)
	items, err := deviceService(c).ListDevices(user.OpenID)
	jsonResult(c, items, err)
}

func BindDevice(c *gin.Context) {
	user := mustMiniappUser(c)
	var req bindDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	item, err := deviceService(c).BindDevice(user.OpenID, req.DeviceID, req.Name)
	jsonResult(c, item, err)
}

func GetDevice(c *gin.Context) {
	user := mustMiniappUser(c)
	item, err := deviceService(c).GetDevice(user.OpenID, c.Param("device_id"))
	jsonResult(c, item, err)
}

func RenameDevice(c *gin.Context) {
	user := mustMiniappUser(c)
	var req renameDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	item, err := deviceService(c).RenameDevice(user.OpenID, c.Param("device_id"), req.Name)
	jsonResult(c, item, err)
}

func DeleteDevice(c *gin.Context) {
	user := mustMiniappUser(c)
	err := deviceService(c).DeleteDevice(user.OpenID, c.Param("device_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "delete device failed", "error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func SendCommand(c *gin.Context) {
	user := mustMiniappUser(c)
	var req commandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	command, event, err := deviceService(c).SendCommand(user.OpenID, c.Param("device_id"), req.Cmd, req.Params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "send command failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg_id": command.MsgID,
		"cmd":    command.Cmd,
		"state":  command.State,
		"detail": command.Detail,
		"event":  event,
	})
}

func SendVoice(c *gin.Context) {
	user := mustMiniappUser(c)
	var req voiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	command, event, err := deviceService(c).SendVoice(user.OpenID, c.Param("device_id"), req.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "send voice failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg_id": command.MsgID,
		"cmd":    command.Cmd,
		"state":  command.State,
		"detail": command.Detail,
		"event":  event,
	})
}

func GetStatus(c *gin.Context) {
	user := mustMiniappUser(c)
	status, err := deviceService(c).GetStatus(user.OpenID, c.Param("device_id"))
	jsonResult(c, status, err)
}

func ListEvents(c *gin.Context) {
	user := mustMiniappUser(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, err := deviceService(c).ListEvents(user.OpenID, c.Param("device_id"), limit)
	jsonResult(c, items, err)
}

func deviceService(c *gin.Context) *omnitouch.Service {
	return omnitouch.NewService(cosy.UseDB(c))
}

func jsonResult(c *gin.Context, data any, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func mustMiniappUser(c *gin.Context) *model.MiniappUser {
	return c.MustGet("miniapp_user").(*model.MiniappUser)
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := internalMiniapp.CurrentUser(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized", "error": err.Error()})
			return
		}
		c.Set("miniapp_user", user)
		c.Next()
	}
}
