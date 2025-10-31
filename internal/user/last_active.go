package user

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/logger"
	cModel "github.com/uozi-tech/cosy/model"
	"github.com/uozi-tech/cosy/redis"
	"github.com/uozi-tech/cosy/settings"
)

type LastActive struct {
	ID         uint64
	LastActive int64
}

var mutex sync.Mutex

// PersistLastActive Persist last active time to redis.
func PersistLastActive(ctx context.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	keys, err := redis.Keys("user:*:last_active")
	if err != nil {
		logger.Error(err)
		return
	}

	defer func() {
		keysLen := len(keys)
		if keysLen > 0 {
			logger.Infof("Persisted user last active, updated %d users", keysLen)
		}
	}()

	var lastActive []LastActive
	for _, key := range keys {
		split := strings.Split(key, ":")
		if len(split) < 2 {
			return
		}
		id := cast.ToUint64(split[1])
		value, err := redis.Get(key)
		if err != nil {
			continue
		}
		lastActive = append(lastActive, LastActive{
			ID:         id,
			LastActive: cast.ToInt64(value),
		})
	}
	if len(lastActive) == 0 {
		return
	}
	err = BatchUpdateLastActive(ctx, lastActive)
	if err != nil {
		logger.Error(err)
	}
	_ = redis.Del(keys...)
}

func BatchUpdateLastActive(ctx context.Context, users []LastActive) error {
	db := cModel.UseDB(ctx)

	createTempTableSQL := `
		CREATE TEMPORARY TABLE IF NOT EXISTS tmp_users (
			id BIGINT PRIMARY KEY,
			last_active BIGINT
		);
	`
	if err := db.Exec(createTempTableSQL).Error; err != nil {
		return fmt.Errorf("failed to create temporary table: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("INSERT INTO tmp_users (id, last_active) VALUES ")
	userLenSub1 := len(users) - 1
	for i, user := range users {
		_, _ = fmt.Fprintf(&sb, "(%d, '%d')", user.ID, user.LastActive)
		if i < userLenSub1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" ON DUPLICATE KEY UPDATE last_active = VALUES(last_active);")

	if err := db.Exec(sb.String()).Error; err != nil {
		return fmt.Errorf("failed to insert into temporary table: %w", err)
	}

	updateSQL := fmt.Sprintf("UPDATE `%susers` u JOIN tmp_users tmp "+
		"ON u.id = tmp.id SET u.last_active = tmp.last_active;", settings.DataBaseSettings.TablePrefix)

	if err := db.Exec(updateSQL).Error; err != nil {
		return fmt.Errorf("failed to update users table: %w", err)
	}

	dropTempTableSQL := "DROP TEMPORARY TABLE IF EXISTS tmp_users;"
	if err := db.Exec(dropTempTableSQL).Error; err != nil {
		return fmt.Errorf("failed to drop temporary table: %w", err)
	}

	return nil
}
