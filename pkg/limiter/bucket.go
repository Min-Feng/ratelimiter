package limiter

import (
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type bucketOption struct {
	maxLimitCount   int32
	interval        time.Duration
	expireDuration  time.Duration
	expireAfterFunc func()
}

func newBucket(o *bucketOption) *bucket {
	bk := &bucket{
		maxLimitCount:   o.maxLimitCount,
		countDuration:   o.interval,
		expireDuration:  o.expireDuration,
		expireAfterFunc: o.expireAfterFunc,
	}

	bk.init()
	return bk
}

type bucket struct {
	// count logic
	maxLimitCount int32
	countDuration time.Duration
	countTicker   *time.Ticker
	count         int32

	// remove logic
	expireDuration    time.Duration
	expireAfterFunc   func()
	expiredTimer      *time.Timer
	resetExpiredTimer chan struct{}
	remove            chan struct{}
	isRemove          int32
}

func (bk *bucket) init() {
	bk.countTicker = time.NewTicker(bk.countDuration)

	bk.expiredTimer = time.AfterFunc(bk.expireDuration, bk.expireAfterFunc)
	bk.resetExpiredTimer = make(chan struct{})

	bk.remove = make(chan struct{})

	go func() {
		for {
			select {
			case <-bk.countTicker.C:
				atomic.StoreInt32(&bk.count, 0)

			case <-bk.expiredTimer.C:
				bk.delete()
			case <-bk.resetExpiredTimer:
				bk.expiredTimer.Reset(bk.expireDuration)

			case <-bk.remove:
				bk.expireAfterFunc()
				bk.countTicker.Stop()
				bk.expiredTimer.Stop()
				return
			}
		}
	}()
}

func (bk *bucket) allow() (count int32, err error) {
	trueValue := int32(1)
	if bk.isRemove == trueValue {
		err = errors.New("bucket channel is closed")
		return
	}

	newCount := atomic.AddInt32(&bk.count, 1)
	if newCount > bk.maxLimitCount {
		return newCount, ErrExceedMaxCount
	}
	bk.resetExpiredTimer <- struct{}{}
	return newCount, nil
}

func (bk *bucket) delete() {
	trueValue := int32(1)
	falseValue := int32(0)
	if atomic.CompareAndSwapInt32(&bk.isRemove, falseValue, trueValue) {
		close(bk.remove)
	}
}
