package limiter

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_bucket_allow(t *testing.T) {
	type result struct {
		userNumber  int32
		returnCount int32
		err         error
	}

	var mu sync.Mutex
	var results = make([]result, 1)
	var wg sync.WaitGroup

	interval := time.Minute
	option := &bucketOption{
		maxLimitCount:   int32(1000),
		interval:        interval,
		expireDuration:  time.Minute,
		expireAfterFunc: func() {},
	}
	bk := newBucket(option)
	concurrencyCount := int(option.maxLimitCount * 2)

	for i := 1; i <= concurrencyCount; i++ {
		wg.Add(1)
		user := i
		go func() {
			defer wg.Done()
			count, err := bk.allow()
			r := result{
				userNumber:  int32(user),
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
		assert.Equalf(t, expectedCount, results[j].returnCount, "userNumber=%v", results[j].userNumber)
		if expectedCount <= option.maxLimitCount {
			assert.NoErrorf(t, results[j].err, "userNumber=%v", results[j].userNumber)
		}
		if expectedCount > option.maxLimitCount {
			assert.Errorf(t, results[j].err, "userNumber=%v", results[j].userNumber)
		}
	}
}
