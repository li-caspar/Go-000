package ratelimit

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRollingChan(t *testing.T) {
	type args struct {
		accuracy      time.Duration
		snippet       time.Duration
		allowRequests int64
	}
	tests := []struct {
		name string
		args args
		want *RollingChan
	}{
		{name: "0", args: struct {
			accuracy      time.Duration
			snippet       time.Duration
			allowRequests int64
		}{accuracy: time.Millisecond * 100, snippet: time.Second, allowRequests: 100}, want: nil},
	}
	var (
		SuccessCount int64
		FailCount    int64
	)
	var f *RollingChan
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f = NewRollingChan(tt.args.snippet, tt.args.accuracy, tt.args.allowRequests)
			for range [10]struct{}{} {
				<-time.After(f.snippet)
				go func() {
					select {
					case <-time.After(f.accuracy):
						//t.Log(time.Now().Second())
						for range [100]struct{}{} {
							if err := f.Take(); err != nil {
								//t.Logf("%d Take() error = %v", i, err)
								atomic.AddInt64(&FailCount, 1)
							} else {
								//t.Log(i, "take once request")
								atomic.AddInt64(&SuccessCount, 1)
							}
						}
					}
				}()
			}
		})
	}
	time.Sleep(2 * time.Second)

	t.Log("SuccessCount:", SuccessCount)
	t.Log("FailCount:", FailCount)
	t.Log("AllCount:", SuccessCount+FailCount)
}
