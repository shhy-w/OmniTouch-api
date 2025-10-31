package user

import (
	"context"
	"encoding/json"
	"testing"

	"git.uozi.org/uozi/cosy-example-api/internal/acl"
	"git.uozi.org/uozi/cosy-example-api/internal/helper"
	"git.uozi.org/uozi/cosy-example-api/model"
	"git.uozi.org/uozi/cosy-example-api/query"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/redis"
	"github.com/uozi-tech/cosy/sandbox"
)

func TestGetByID(t *testing.T) {
	sandbox.NewInstance("../../app.testing.ini", "mysql").
		RegisterModels(model.User{}, model.UserGroup{}).
		Run(func(instance *sandbox.Instance) {
			db := cosy.UseDB(context.Background())
			query.Init(db)
			model.Use(db)
			// create user group
			group := &model.UserGroup{
				Name: "test",
				Permissions: acl.PermissionList{
					{
						Subject: acl.User,
						Action:  acl.Write,
					},
				},
			}

			err := db.Create(group).Error
			if err != nil {
				t.Error(err)
			}
			// create user
			user := &model.User{
				Name:        "test",
				Password:    "test",
				Email:       "test",
				UserGroupID: group.ID,
			}

			err = db.Create(user).Error
			if err != nil {
				t.Error(err)
			}

			user, err = GetByID(user.ID)
			if err != nil {
				t.Error(err)
			}
			// check cache
			key := helper.BuildUserKey(user.ID)
			value, err := redis.Get(key)
			if err != nil {
				t.Error(err)
			}
			cachedUser := &model.User{}
			err = json.Unmarshal([]byte(cast.ToString(value)), cachedUser)
			if err != nil {
				t.Error(err)
			}
			assert.ObjectsAreEqual(user, cachedUser)
			assert.ObjectsAreEqual(user.UserGroup, group)
		})
}
