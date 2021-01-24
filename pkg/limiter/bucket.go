package limiter

import (
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// maxLimitCount is number that allowed per interval.
// interval reset time
func newBucket(maxLimitCount int32, interval time.Duration) *bucket {
	bk := &bucket{
		maxLimitCount: maxLimitCount,
		interval:      interval,
	}

	bk.initResetLimitInterval()
	return bk
}

type bucket struct {
	maxLimitCount int32
	interval      time.Duration

	ticker *time.Ticker
	stop   chan struct{}
	isStop int32

	count int32
}

func (bk *bucket) initResetLimitInterval() {
	stop := make(chan struct{})
	ticker := time.NewTicker(bk.interval)

	bk.stop = stop
	bk.ticker = ticker

	go func() {
		for {
			select {
			case <-stop:
				ticker.Stop()
				return
			case <-ticker.C:
				atomic.StoreInt32(&bk.count, 0)
			}
		}
	}()
}

func (bk *bucket) allow() (count int32, err error) {
	trueValue := int32(1)
	if bk.isStop == trueValue {
		err = errors.New("bucket channel is closed")
		return
	}

	newCount := atomic.AddInt32(&bk.count, 1)
	if newCount > bk.maxLimitCount {
		return newCount, ErrExceedMaxCount
	}
	return newCount, nil
}

func (bk *bucket) close() {
	trueValue := int32(1)
	falseValue := int32(0)
	if atomic.CompareAndSwapInt32(&bk.isStop, falseValue, trueValue) {
		close(bk.stop)
	}
	return
}
