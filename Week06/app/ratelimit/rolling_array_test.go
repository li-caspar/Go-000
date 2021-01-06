package ratelimit

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewRollingArray(t *testing.T) {
	ra := NewRollingArray(time.Second, time.Millisecond*100, 100)
	rand.Seed(time.Now().UnixNano())
	//for i := 0; i < 5; i++ {
	//time.Sleep(time.Second)
	//go func() {
	//time.Sleep(time.Millisecond * 10)
	/*for j := 0; j < 100; j++ {
		err := ra.Take(i)
		if err == nil {
			//fmt.Printf("success row:%d\n", j)
		} else {
			//fmt.Printf("fail row:%d\n", j)
		}
	}*/
	//}()
	//time.Sleep(time.Millisecond*1000 + 100)
	//}
	count := 0
	failCount := 0
	for i := range [10]struct{}{} {
		<-time.After(time.Second)
		go func(i int) {
			select {
			case <-time.After(time.Millisecond * 100):
				for range [50]struct{}{} {
					if err := ra.Take(); err != nil {
						//t.Errorf("Take() error = %v", err)
						failCount++
					} else {
						//t.Log(i, "take once request")
						count++
					}
				}
			}
		}(i)
	}

	time.Sleep(2 * time.Second)
	if count != 100*5 {
		t.Error("count:", count)
	}
	if failCount > 0 {
		t.Fatalf("failCount:%d", failCount)
	}
	//fmt.Printf("%v", ra)
}

func TestRound(t *testing.T) {
	for i := 0; i <= 10; i++ {
		time.Sleep(time.Second)
		fmt.Printf("%d\n", time.Now().UnixNano()/1e6)
	}
}
