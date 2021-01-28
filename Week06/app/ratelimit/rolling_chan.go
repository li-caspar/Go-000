package ratelimit

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var onceChan sync.Once

type RollingChan struct {
	values          chan int64 //用来储存每个格子的请求数
	size            int
	currentRequests int64         //当前格子的请求数
	allRequests     int64         //窗口的所有请求数
	allowRequests   int64         //允许的最大的请求数
	snippet         time.Duration //窗口的长度(时间为单位)
	accuracy        time.Duration //每个格的长度(时间为单位)
	log             strings.Builder
}

var _ RateLimiter = &RollingChan{}

func NewRollingChan(snippet time.Duration, accuracy time.Duration, allowRequests int64) *RollingChan {
	size := int(snippet / accuracy) //应该有多少个格子
	if size < 0 {
		size = 0
	}
	r := &RollingChan{
		values:        make(chan int64, size),
		size:          size,
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

func (r *RollingChan) writeLog(str string) {
	r.log.WriteString(str)
}

func (r *RollingChan) GetString() string {
	return r.log.String()
}

//滑入
func slideIn(r *RollingChan) {
	for {
		select {
		case <-time.After(r.accuracy):
			requests := atomic.SwapInt64(&r.currentRequests, 0) //把当前格子的请求数清空
			r.values <- requests
			r.writeLog(fmt.Sprintf("time:%d:IN t:%d, currentRequests:%d, allRequests:%d\n", time.Now().UnixNano()/1e6, requests, r.currentRequests, r.allRequests))
		}
	}
}

//滑出
func slideOut(r *RollingChan) {
	for {
		select {
		case <-time.After(r.accuracy):
			if len(r.values) != r.size {
				continue
			}
			t := <-r.values //从chan中移出一个格子
			r.writeLog(fmt.Sprintf("time:%d:OUT t:%d currentRequests:%d, allRequests:%d\n", time.Now().UnixNano()/1e6, t, r.currentRequests, r.allRequests))
			if t != 0 {
				atomic.AddInt64(&r.allRequests, -t) //将窗口的请求数减少
			}
		}
	}
}
