package limiter

import (
	"sync"
	"time"
)

func NewLocalLimiter(maxLimitCount int32, interval time.Duration) *LocalLimiter {
	return &LocalLimiter{
		maxLimitCount:      maxLimitCount,
		resetCountInterval: interval,
		keyExpireDuration:  2 * interval,
	}
}

type LocalLimiter struct {
	maxLimitCount      int32
	resetCountInterval time.Duration
	keyCollection      sync.Map
	keyExpireDuration  time.Duration
}

type getBucket func() *bucket

func (l *LocalLimiter) Allow(key string) (count int32, err error) {
	fn, ok := l.keyCollection.Load(key)
	if ok {
		oldBucket := fn.(getBucket)()
		return oldBucket.allow()
	}

	var freshBucket *bucket
	var once sync.Once

	lazyInit := func() *bucket {
		once.Do(func() {
			option := &bucketOption{
				maxLimitCount:  l.maxLimitCount,
				interval:       l.resetCountInterval,
				expireDuration: l.keyExpireDuration,
				expireAfterFunc: func() {
					l.keyCollection.Delete(key)
				},
			}
			freshBucket = newBucket(option)
		})
		return freshBucket
	}

	fn, loaded := l.keyCollection.LoadOrStore(key, getBucket(lazyInit))
	if !loaded {
		freshBucket = lazyInit()

		// 減少之後 load 時, once 造成的效能差異
		l.keyCollection.Store(key, getBucket(func() *bucket { return freshBucket }))
		return freshBucket.allow()
	}

	oldBucket := fn.(getBucket)() // 執行從 sync.Map 拿到的函數
	return oldBucket.allow()
}

func (l *LocalLimiter) Delete(key string) error {
	fn, ok := l.keyCollection.Load(key)
	if !ok {
		return nil
	}

	l.keyCollection.Delete(key)
	oldBucket := fn.(getBucket)()
	oldBucket.delete()
	return nil
}
