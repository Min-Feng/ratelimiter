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
	wg            sync.WaitGroup
}

func (l *LocalLimiter) Allow(key string) (count int32, err error) {
	l.wg.Wait()
	v, ok := l.keyCollection.Load(key)
	if ok {
		return v.(*bucket).allow()
	}

	var bk *bucket
	l.wg.Add(1)

	_, loaded := l.keyCollection.LoadOrStore(key, bk)
	if !loaded {
		bk = newBucket(l.maxLimitCount, l.interval)
		l.keyCollection.Store(key, bk)
		l.wg.Done()
		return bk.allow()
	}

	l.wg.Done()
	l.wg.Wait()
	v, _ = l.keyCollection.Load(key)
	return v.(*bucket).allow()
}

func (l *LocalLimiter) Close(key string) error {
	panic("implement me")
}

func (l *LocalLimiter) CloseAll() error {
	panic("implement me")
}
