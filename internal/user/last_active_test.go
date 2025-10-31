package user

// import (
// 	"testing"

// 	"git.uozi.org/uozi/hotel-api/model"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/uozi-tech/cosy"
// 	"github.com/uozi-tech/cosy/sandbox"
// )

// func TestBatchUpdateLastActive(t *testing.T) {
// 	sandbox.NewInstance("../../app.testing.ini", "mysql").
// 		RegisterModels(model.User{}, model.UserGroup{}).
// 		Run(func(instance *sandbox.Instance) {
// 			// create user
// 			db := cosy.UseDB()
// 			users := []model.User{
// 				{
// 					Name:     "test",
// 					Password: "test",
// 					Email:    "test",
// 				},
// 				{
// 					Name:     "test1",
// 					Password: "test1",
// 					Email:    "test1",
// 				},
// 			}
// 			err := db.Create(users).Error
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			var lastActive []LastActive
// 			for _, u := range users {
// 				active := u.UpdateLastActive()
// 				lastActive = append(lastActive, LastActive{
// 					ID:         u.ID,
// 					LastActive: active,
// 				})
// 			}
// 			// batch update last active
// 			err = BatchUpdateLastActive(lastActive)

// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}

// 			// check last active
// 			u := &model.User{}
// 			err = db.First(u, 1).Error
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			assert.Equal(t, lastActive[0].LastActive, u.LastActive)
// 			u = &model.User{}
// 			err = db.First(u, 2).Error
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			assert.Equal(t, lastActive[1].LastActive, u.LastActive)
// 		})
// }

// func TestPersistLastActive(t *testing.T) {
// 	sandbox.NewInstance("../../app.testing.ini", "mysql").
// 		RegisterModels(model.User{}, model.UserGroup{}).
// 		Run(func(instance *sandbox.Instance) {
// 			// create user
// 			db := cosy.UseDB()
// 			users := []model.User{
// 				{
// 					Name:     "test",
// 					Password: "test",
// 					Email:    "test",
// 					Phone:    "test",
// 				},
// 				{
// 					Name:     "test1",
// 					Password: "test1",
// 					Email:    "test1",
// 					Phone:    "test1",
// 				},
// 			}
// 			err := db.Create(users).Error
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			var lastActive []LastActive
// 			for _, u := range users {
// 				lastActive = append(lastActive, LastActive{
// 					ID:         u.ID,
// 					LastActive: u.UpdateLastActive(),
// 				})
// 			}

// 			// Test at here
// 			PersistLastActive()

// 			// check last active
// 			u := &model.User{}
// 			err = db.First(u, 1).Error
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			assert.Equal(t, lastActive[0].LastActive, u.LastActive)
// 			u = &model.User{}
// 			err = db.First(u, 2).Error
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			assert.Equal(t, lastActive[1].LastActive, u.LastActive)
// 		})

// }
