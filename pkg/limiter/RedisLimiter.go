package limiter

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
)

func NewRedisLimiter(client *redis.Client, maxLimitCount int32, interval time.Duration) *RedisLimiter {
	type redisKey = string
	limiter := &RedisLimiter{
		maxLimitCount:      int64(maxLimitCount),
		resetCountInterval: interval,
		enableResetCount:   make(chan redisKey),
		closeResetCount:    make(chan redisKey),
		allCloseNotify:     make(chan struct{}),

		client:            client,
		prefixKey:         "limiter:",
		keyExpireDuration: 2 * interval,
	}
	limiter.resetCountManager()
	return limiter
}

type RedisLimiter struct {
	maxLimitCount      int64
	resetCountInterval time.Duration
	enableResetCount   chan string
	closeResetCount    chan string
	allCloseNotify     chan struct{}

	client            *redis.Client
	prefixKey         string
	keyExpireDuration time.Duration
}

func (r *RedisLimiter) Allow(key string) (int32, error) {
	redisKey := r.redisKey(key)

	count, err := r.client.Incr(redisKey).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis command incr by redis key=%v", redisKey)
	}

	if count > r.maxLimitCount {
		return int32(count), ErrExceedMaxCount
	}

	if r.firstTime(count) {
		_, err := r.client.Expire(redisKey, r.keyExpireDuration).Result()
		if err != nil {
			return 0, errors.Wrapf(err, "redis set expire time by key=%v ", redisKey)
		}
		r.enableResetCount <- redisKey
	}

	return int32(count), nil
}

func (r *RedisLimiter) resetCountManager() {
	type redisKey = string
	type keyCloseNotify = chan struct{}
	manager := make(map[redisKey]keyCloseNotify)

	go func() {
		for {
			select {
			case k := <-r.enableResetCount:
				if _, ok := manager[k]; ok {
					break
				}
				notify := r.createResetCountWorker(k)
				manager[k] = notify

			case k := <-r.closeResetCount:
				notify, ok := manager[k]
				if ok {
					break
				}
				delete(manager, k)
				close(notify)

			case <-r.allCloseNotify:
				for key, keyNotify := range manager {
					delete(manager, key)
					close(keyNotify)
				}
				return
			}
		}
	}()
}

func (r *RedisLimiter) createResetCountWorker(redisKey string) chan struct{} {
	ticker := time.NewTicker(r.resetCountInterval)
	keyCloseNotify := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_, err := r.client.Set(redisKey, 0, r.keyExpireDuration).Result()
				if err != nil {
					panic(err)
				}
			case <-keyCloseNotify:
				ticker.Stop()
				return
			}
		}
	}()
	return keyCloseNotify
}

func (r *RedisLimiter) redisKey(key string) string {
	return r.prefixKey + key
}

func (r *RedisLimiter) firstTime(count int64) bool {
	return count == 1
}

func (r *RedisLimiter) Delete(key string) error {
	redisKey := r.redisKey(key)
	err := r.client.Del(redisKey).Err()
	if err != nil {
		return errors.Wrapf(err, "redis delete by key=%v", redisKey)
	}
	r.closeResetCount <- redisKey
	return nil
}
