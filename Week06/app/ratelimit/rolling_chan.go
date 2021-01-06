package ratelimit

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var onceChan sync.Once

type RollingChan struct {
	values          chan int64    //用来储存每个格子的请求数
	currentRequests int64         //当前格子的请求数
	allRequests     int64         //窗口的所有请求数
	allowRequests   int64         //允许的最大的请求数
	snippet         time.Duration //窗口的长度(时间为单位)
	accuracy        time.Duration //每个格的长度(时间为单位)
}

var _ RateLimiter = &RollingChan{}

func NewRollingChan(snippet time.Duration, accuracy time.Duration, allowRequests int64) *RollingChan {
	size := snippet/accuracy - 1 //应该有多少个格子
	if size < 0 {
		size = 0
	}
	r := &RollingChan{
		values:        make(chan int64, size),
		allowRequests: allowRequests,
		snippet:       snippet,
		accuracy:      accuracy,
	}
	return r
}

func (r *RollingChan) Take() error {
	onceChan.Do(func() {
		go func() {
			if err := recover(); err != nil {
				fmt.Printf("slideOut error:%s", err)
			}
			slideOut(r)
		}()
		go func() {
			if err := recover(); err != nil {
				fmt.Printf("slideIn error:%s", err)
			}
			slideIn(r)
		}()

	})
	if atomic.LoadInt64(&r.allRequests) >= r.allowRequests {
		return ErrExceededLimit
	}
	atomic.AddInt64(&r.currentRequests, 1)
	atomic.AddInt64(&r.allRequests, 1)
	return nil
}

//滑入
func slideIn(r *RollingChan) {
	for {
		select {
		case <-time.After(r.accuracy):
			requests := atomic.SwapInt64(&r.currentRequests, 0) //把当前格子的请求数清空
			r.values <- requests
		}
	}
}

//滑出
func slideOut(r *RollingChan) {
	for {
		select {
		case <-time.After(r.accuracy):
			t := <-r.values                     //从chan中移出一个格子
			atomic.AddInt64(&r.allRequests, -t) //将窗口的请求数减少
		}
	}
}
