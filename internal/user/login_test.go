package user

import (
	"context"
	"os"
	"testing"

	"git.uozi.org/uozi/cosy-example-api/model"
	"git.uozi.org/uozi/cosy-example-api/query"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/sandbox"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	if os.Getenv("OMNITOUCH_RUN_DB_TESTS") != "1" {
		t.Skip("skip legacy database integration test")
	}
	sandbox.NewInstance("../../app.testing.ini", "mysql").
		RegisterModels(model.User{}).
		Run(func(instance *sandbox.Instance) {
			db := cosy.UseDB(context.Background())
			query.Init(db)
			model.Use(db)

			pwd, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

			db.Create(&[]model.User{{
				Email:    "test",
				Password: string(pwd),
				Status:   model.UserStatusActive,
			}, {
				Email:    "banned",
				Password: string(pwd),
				Status:   model.UserStatusBan,
			}})

			t.Run("success", func(t *testing.T) {
				user, err := Login(&gin.Context{}, "test", "test")
				if err != nil {
					t.Error(err)
					return
				}
				assert.NotNil(t, user)
			})

			t.Run("password incorrect", func(t *testing.T) {
				_, err := Login(&gin.Context{}, "test", "123456")
				assert.Equal(t, ErrPasswordIncorrect, err)
				_, err = Login(&gin.Context{}, "123456", "123456")
				assert.Equal(t, ErrPasswordIncorrect, err)
			})

			t.Run("user banned", func(t *testing.T) {
				_, err := Login(&gin.Context{}, "banned", "test")
				assert.Equal(t, ErrUserBanned, err)
			})
		})
}
