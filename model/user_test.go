package model

import (
	"context"
	"testing"

	"git.uozi.org/uozi/cosy-example-api/internal/acl"
	"git.uozi.org/uozi/cosy-example-api/internal/helper"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/uozi-tech/cosy"
	cModel "github.com/uozi-tech/cosy/model"
	"github.com/uozi-tech/cosy/redis"
	"github.com/uozi-tech/cosy/sandbox"
)

func TestModelUserGetPermissionsMap(t *testing.T) {
	sandbox.NewInstance("../app.testing.ini", "mysql").
		RegisterModels(User{}, UserGroup{}).
		Run(func(instance *sandbox.Instance) {
			Use(cModel.UseDB(context.Background()))
			// create a user group
			group := &UserGroup{
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
			user := &User{
				Name:        "test",
				Password:    "test",
				Email:       "test",
				UserGroupID: group.ID,
			}
			assert.Equal(t, uint64(1), group.ID)
			err = db.Create(user).Error
			if err != nil {
				t.Error(err)
			}
			group.CleanCache()
			permissionMap := user.GetPermissionsMap()
			t.Log(permissionMap)
			assert.Equal(t, true, permissionMap[acl.User][acl.Write])
			assert.Equal(t, true, acl.Can(permissionMap, acl.User, acl.Write))
			assert.Equal(t, true, acl.Can(permissionMap, acl.User, acl.Read))

			group.CleanCache()

			value, _ := redis.Get(getUserGroupKey(group.ID))
			assert.Equal(t, "", value)
		})
}

func TestUpdateLastActive(t *testing.T) {
	sandbox.NewInstance("../app.testing.ini", "mysql").
		RegisterModels(User{}, UserGroup{}).
		Run(func(instance *sandbox.Instance) {
			Use(cosy.UseDB(context.Background()))
			user := &User{
				Model: Model{
					ID: 1,
				},
			}

			now := user.UpdateLastActive()

			key := helper.BuildUserKey(1, "last_active")

			value, err := redis.Get(key)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, now, cast.ToInt64(value))
		})
}

func TestUserIsAdmin(t *testing.T) {
	sandbox.NewInstance("../app.testing.ini", "mysql").
		RegisterModels(User{}, UserGroup{}).
		Run(func(instance *sandbox.Instance) {
			Use(cosy.UseDB(context.Background()))

			// create an admin group
			group := &UserGroup{
				Name: "admin",
				Permissions: acl.PermissionList{
					{
						Subject: acl.All,
						Action:  acl.Write,
					},
				},
			}
			err := db.Create(group).Error
			if err != nil {
				t.Error(err)
			}
			// create user group
			uGroup := &UserGroup{
				Name: "user",
				Permissions: acl.PermissionList{
					{
						Subject: acl.User,
						Action:  acl.Read,
					},
				},
			}
			err = db.Create(uGroup).Error
			if err != nil {
				t.Error(err)
			}
			// create admin
			admin := &User{
				Name:        "admin",
				Password:    "admin",
				Email:       "test",
				UserGroupID: group.ID,
			}

			err = db.Create(admin).Error
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, true, admin.IsAdmin())

			user := &User{
				Name:        "test",
				Password:    "test",
				Email:       "test1",
				UserGroupID: uGroup.ID,
			}

			err = db.Create(user).Error
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, false, user.IsAdmin())
		})
}
