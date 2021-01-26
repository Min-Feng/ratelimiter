package limiter

import (
	"github.com/Min-Feng/ratelimiter/pkg/configs"
	"github.com/Min-Feng/ratelimiter/pkg/infra"
)

type Limiter interface {
	Allow(key string) (count int32, err error)
	Delete(key string) error
}

func New(cfg *configs.Config, kind string) Limiter {
	var l Limiter
	switch kind {
	case "redis":
		client := infra.NewRedis(&cfg.Redis)
		l = NewRedisLimiter(client, cfg.Limiter.MaxLimitCount, cfg.Limiter.ResetCountInterval())
	case "local":
		l = NewLocalLimiter(cfg.Limiter.MaxLimitCount, cfg.Limiter.ResetCountInterval())
	default:
		panic("not support kind=" + kind)
	}
	return l
}
