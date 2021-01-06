package ratelimit

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var onceArr sync.Once

type RollingArray struct {
	values        []int64 //环形slice， 大小由snippet/accuracy决定, 数值为每个格的请求数
	size          int
	mutex         sync.RWMutex
	lastIndex     int32
	snippet       time.Duration //窗口的长度(时间为单位)
	accuracy      time.Duration //每个格的长度(时间为单位)
	allowRequests int64         //窗口允许最大的请求数
}

var _ RateLimiter = &RollingArray{}

func NewRollingArray(snippet time.Duration, accuracy time.Duration, allowRequests int64) *RollingArray {
	var size int = int(snippet) / int(accuracy)
	return &RollingArray{
		values:        make([]int64, size, size),
		size:          size,
		lastIndex:     0,
		snippet:       snippet,
		accuracy:      accuracy,
		allowRequests: allowRequests,
	}
}

//是否允许通过请求
func (r *RollingArray) Take() error {
	sum := r.sum()
	if sum > r.allowRequests {
		return ErrExceededLimit
	}
	onceArr.Do(func() {
		go func() {
			if err := recover(); err != nil {
				fmt.Printf("recover error:%s", err)
			}
			for {
				select {
				case <-time.After(r.accuracy):
					r.slide()
				}
			}
		}()
		time.Sleep(r.accuracy)
	})
	r.increment()
	return nil
}

//获取窗口的最新总请求数
func (r *RollingArray) sum() int64 {
	var sum int64
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, value := range r.values {
		sum += value
	}
	return sum
}

func (r *RollingArray) string() string {
	var str strings.Builder
	for i, value := range r.values {
		str.WriteString(fmt.Sprintf("%d:%d  ", i, value))
	}
	return str.String()
}

//最新格的请求数自增
func (r *RollingArray) increment() {
	index := r.getLastIndex()
	value := atomic.LoadInt64(&r.values[index])
	atomic.CompareAndSwapInt64(&r.values[index], value, value+1)
}

//获取最新格的索引
func (r *RollingArray) getLastIndex() int32 {
	return atomic.LoadInt32(&r.lastIndex)
}

//滑动窗口
func (r *RollingArray) slide() {
	index := r.getLastIndex()
	index = (index + 1) % int32(r.size)
	atomic.StoreInt64(&r.values[index], 0)
	atomic.StoreInt32(&r.lastIndex, index)
}
