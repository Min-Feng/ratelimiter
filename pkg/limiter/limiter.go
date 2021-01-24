package limiter

import (
	"sync"
	"time"
)

type Limiter interface {
	Allow(key string) (count int32, err error)
	Delete(key string) error
	DeleteAll() error
}

func NewLocalLimiter(maxLimitCount int32, interval time.Duration) *LocalLimiter {
	return &LocalLimiter{
		maxLimitCount:     maxLimitCount,
		interval:          interval,
		keyExpireDuration: time.Hour,
	}
}

type LocalLimiter struct {
	maxLimitCount int32
	interval      time.Duration

	keyCollection     sync.Map
	keyExpireDuration time.Duration
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

func (l *LocalLimiter) Delete(key string) error {
	fn, ok := l.keyCollection.Load(key)
	if !ok {
		return nil
	}
	l.delete(key, fn)
	return nil
}

func (l *LocalLimiter) DeleteAll() error {
	l.keyCollection.Range(func(key, value interface{}) bool {
		l.delete(key.(string), value)
		return true
	})
	return nil
}

func (l *LocalLimiter) delete(key string, fn interface{}) {
	type Getter func() *bucket
	l.keyCollection.Delete(key)
	oldBucket := fn.(Getter)()
	oldBucket.close()
}
