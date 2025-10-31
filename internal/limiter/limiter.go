package limiter

import (
	"git.uozi.com/uozi/rate-limiter-go"
	"github.com/uozi-tech/cosy/redis"
)

var limiter *rate_limiter.Limiter

func Init() {
	limiter = rate_limiter.NewLimiter(redis.GetClient())
}

func GetLimiter() *rate_limiter.Limiter {
	return limiter
}
