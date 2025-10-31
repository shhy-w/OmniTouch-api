package user

import (
	"git.uozi.org/uozi/cosy-example-api/model"
	"git.uozi.org/uozi/cosy-example-api/query"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/logger"
	"golang.org/x/crypto/bcrypt"
	// "git.uozi.org/uozi/cosy-example-api/api/global"
)

var (
	e                        = cosy.NewErrorScope("user")
	ErrPasswordIncorrect     = e.New(4031, "password incorrect")
	ErrMaxAttempts           = e.New(4291, "max attempts exceeded")
	ErrUserBanned            = e.New(4033, "user banned")
	ErrUserWaitForValidation = e.New(4034, "user wait for validation")
)

func Login(c *gin.Context, email string, password string) (user *model.User, err error) {
	u := query.User

	user, err = u.Where(u.Email.Eq(email)).Preload(u.UserGroup).First()
	if err != nil {
		logger.Error(err)
		return nil, ErrPasswordIncorrect
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Error(err)
		return nil, ErrPasswordIncorrect
	}

	if user.Status == model.UserStatusBan {
		return nil, ErrUserBanned
	}

	return
}
