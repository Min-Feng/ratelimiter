package limiter

import (
	"sync"
	"time"
)

type Limiter interface {
	Allow(key string) (count int32, err error)
	Close(key string) error
	CloseAll() error
}

func NewLocalLimiter(maxLimitCount int32, interval time.Duration) *LocalLimiter {
	return &LocalLimiter{maxLimitCount: maxLimitCount, interval: interval}
}

type LocalLimiter struct {
	maxLimitCount int32
	interval      time.Duration

	keyCollection sync.Map
}

func (l *LocalLimiter) Allow(key string) (count int32, err error) {
	type Getter func() *bucket

	fn, ok := l.keyCollection.Load(key)
	if ok {
		oldBucket := fn.(Getter)()
		return oldBucket.allow()
	}

	var freshBucket *bucket
	var once sync.Once
	lazyInit := func() *bucket {
		once.Do(func() {
			freshBucket = newBucket(l.maxLimitCount, l.interval)
		})
		return freshBucket
	}

	fn, loaded := l.keyCollection.LoadOrStore(key, Getter(lazyInit))
	if !loaded {
		freshBucket = lazyInit()

		// 減少之後 load 時, once 造成的效能差異
		l.keyCollection.Store(key, Getter(func() *bucket { return freshBucket }))
		return freshBucket.allow()
	}

	oldBucket := fn.(Getter)() // 已經 loaded 的 函數 Allow(key), 執行從 sync.Map 拿到的函數
	return oldBucket.allow()
}

func (l *LocalLimiter) Close(key string) error {
	panic("implement me")
}

func (l *LocalLimiter) CloseAll() error {
	panic("implement me")
}
