// +build integration

package limiter

import (
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/ratelimiter/pkg/configs"
)

func Test_LocalLimiter_Allow_count_successful(t *testing.T) {
	cfg := configs.New("config")
	cfg.Limiter.MaxLimitCount = 1000
	cfg.Limiter.ResetCountIntervalSeconds = 2

	limiter := New(&cfg, "local")
	limiterTestTemplate(t, cfg, limiter)
}

func Test_RedisLimiter_Allow_count_successful(t *testing.T) {
	cfg := configs.New("config")
	cfg.Limiter.MaxLimitCount = 1000
	cfg.Limiter.ResetCountIntervalSeconds = 2

	limiter := New(&cfg, "redis")
	limiterTestTemplate(t, cfg, limiter)
}

func limiterTestTemplate(t *testing.T, cfg configs.Config, limiter Limiter) {
	type result struct {
		returnCount int32
		err         error
	}

	var mu sync.Mutex
	var results = make([]result, 1)
	var wg sync.WaitGroup

	maxLimitCount := cfg.Limiter.MaxLimitCount
	concurrencyCount := int(maxLimitCount * 2)
	key := "192.0.2.1"

	for i := 0; i < concurrencyCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count, err := limiter.Allow(key)
			r := result{
				returnCount: count,
				err:         err,
			}
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}()
	}
	wg.Wait()

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].returnCount < results[j].returnCount
	})

	for j := 1; j <= concurrencyCount; j++ {
		expectedCount := int32(j)
		assert.Equal(t, expectedCount, results[j].returnCount)
		if expectedCount <= maxLimitCount {
			assert.NoError(t, results[j].err)
		}
		if expectedCount > maxLimitCount {
			assert.Error(t, results[j].err)
		}
	}

	err := limiter.Delete(key)
	assert.NoError(t, err)
}
