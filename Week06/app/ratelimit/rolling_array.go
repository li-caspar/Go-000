package ratelimit

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var onceArr sync.Once

type RollingArray struct {
	values        []int64 //环形slice， 大小由snippet/accuracy决定, 数值为每个格的请求数
	size          int     //slice的大小
	mutex         sync.RWMutex
	lastIndex     int32         //最新的格子位置
	snippet       time.Duration //窗口的长度(时间为单位)
	accuracy      time.Duration //每个格的长度(时间为单位)
	allowRequests int64         //窗口允许最大的请求数
}

var _ RateLimiter = &RollingArray{}

func NewRollingArray(snippet time.Duration, accuracy time.Duration, allowRequests int64) *RollingArray {
	size := int(snippet / accuracy)
	r := &RollingArray{
		values:        make([]int64, size, size),
		size:          size,
		lastIndex:     0,
		snippet:       snippet,
		accuracy:      accuracy,
		allowRequests: allowRequests,
	}
	return r
}

//是否允许通过请求
func (r *RollingArray) Take() error {
	onceArr.Do(func() {
		go func() {
			if err := recover(); err != nil {
				fmt.Printf("recover error:%s", err)
			}
			slide(r)
		}()
	})

	sum := r.sum()
	if sum >= r.allowRequests {
		return ErrExceededLimit
	}
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
	r.mutex.Lock()
	r.values[r.lastIndex]++
	defer r.mutex.Unlock()
}

//滑动窗口
func slide(r *RollingArray) {
	for {
		select {
		case <-time.After(r.accuracy):
			r.mutex.Lock()
			index := (r.lastIndex + 1) % int32(r.size) //计算最新格是哪个
			r.values[index] = 0
			r.lastIndex = index
			r.mutex.Unlock()
		}
	}

}
