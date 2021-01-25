package limiter

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocalLimiter_Allow_count_successful(t *testing.T) {
	type result struct {
		returnCount int32
		err         error
	}

	var mu sync.Mutex
	var results = make([]result, 1)
	var wg sync.WaitGroup

	maxLimitCount := int32(1000)
	interval := time.Minute
	limiter := NewLocalLimiter(maxLimitCount, interval)
	concurrencyCount := int(maxLimitCount * 2)

	for i := 0; i < concurrencyCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count, err := limiter.Allow("ip1")
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
}
