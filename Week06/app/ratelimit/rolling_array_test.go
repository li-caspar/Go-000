package ratelimit

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

//TODO
func TestNewRollingArray(t *testing.T) {
	ra := NewRollingArray(time.Second, time.Millisecond*100, 100)
	rand.Seed(time.Now().UnixNano())
	var (
		SuccessCount int64
		FailCount    int64
	)
	for range [10]struct{}{} {
		<-time.After(time.Second)
		for range [100]struct{}{} {
			if err := ra.Take(); err != nil {
				//t.Errorf("Take() error = %v", err)
				atomic.AddInt64(&FailCount, 1)
			} else {
				//t.Log(i, "take once request")
				atomic.AddInt64(&SuccessCount, 1)
			}
		}
		/*go func(i int) {
			select {
			case <-time.After(time.Millisecond * 100):
				for range [101]struct{}{} {
					if err := ra.Take(); err != nil {
						//t.Errorf("Take() error = %v", err)
						atomic.AddInt64(&FailCount, 1)
					} else {
						//t.Log(i, "take once request")
						atomic.AddInt64(&SuccessCount, 1)
					}
				}
			}
		}(i)*/
	}

	time.Sleep(2 * time.Second)

	t.Log("SuccessCount:", SuccessCount)
	t.Log("FailCount:", FailCount)
	t.Log("AllCount:", SuccessCount+FailCount)
}
