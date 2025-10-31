package user

import (
	"git.uozi.org/uozi/cosy-example-api/model"
	"github.com/uozi-tech/cosy"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(ctx *cosy.Ctx[model.User]) {
	if ctx.Payload["password"] == nil {
		return
	}
	pwd := ctx.Payload["password"].(string)
	if pwd != "" {
		pwdBytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			ctx.AbortWithError(err)
			return
		}
		ctx.Model.Password = string(pwdBytes)
	} else {
		delete(ctx.Payload, "password")
	}
}
