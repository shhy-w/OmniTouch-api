package runtime_settings

import (
	"encoding/json"
	"github.com/uozi-tech/cosy/logger"
)

type Auth struct {
	BanThresholdMinutes int `json:"ban_threshold_minutes" binding:"min=1"`
	MaxAttempts         int `json:"max_attempts" binding:"min=1"`
}

func GetAuthSettings() (a *Auth) {
	a = &Auth{
		BanThresholdMinutes: 10,
		MaxAttempts:         10,
	}
	err := Get("auth", a)
	if err != nil {
		logger.Error(err)
	}
	return
}

func UpdateAuthSettings(target *Auth) {
	bytes, err := json.Marshal(target)
	if err != nil {
		logger.Error(err)
		return
	}
	Set("auth", string(bytes))
}
