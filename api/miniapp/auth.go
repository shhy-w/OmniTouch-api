package miniapp

import (
	"net/http"

	internalMiniapp "git.uozi.org/uozi/cosy-example-api/internal/miniapp"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

type wxLoginRequest struct {
	Code string `json:"code"`
}

func WxLogin(c *gin.Context) {
	var req wxLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}
	user, err := internalMiniapp.LoginByCode(cosy.UseDB(c), req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "wechat login failed", "error": err.Error()})
		return
	}
	token, err := internalMiniapp.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "generate token failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":  token,
		"openid": user.OpenID,
	})
}
