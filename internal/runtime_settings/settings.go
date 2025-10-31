package runtime_settings

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/redis"
)

func buildKey(key interface{}) string {
	var sb strings.Builder
	sb.WriteString("settings:")
	sb.WriteString(cast.ToString(key))
	return sb.String()
}

func Get(key string, value interface{}) (err error) {
	key = buildKey(key)
	// key==settings:auth
	logger.Debug("redis_key:" + key)
	raw, err := redis.Get(key)
	if err != nil {
		return
	}
	return json.Unmarshal([]byte(raw), value)
}

func Set(key string, value interface{}) {
	key = buildKey(key)
	_ = redis.Set(key, value, 0)
}
