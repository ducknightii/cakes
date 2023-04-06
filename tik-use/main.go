package main

import (
	"fmt"
	"github.com/andy2046/tik"
	"sync"
	"time"
)

func main() {
	var l sync.RWMutex
	// init a new instance
	tk := tik.New()
	i := 0
	cb := func() {
		l.Lock()
		i++
		fmt.Println("do onece...")
		l.Unlock()
	}
	// schedule to run cb in 500ms
	to := tk.Schedule(500, cb)

	if !to.Pending() {
		panic("it should be pending")
	}

	if to.Expired() {
		panic("it should NOT be expired")
	}

	/*for {
		time.Sleep(100 * time.Millisecond)

		if tk.AnyPending() {
			continue
		}

		if tk.AnyExpired() {
			continue
		}

		break
	}*/

	time.Sleep(2 * time.Second)

	l.RLock()
	defer l.RUnlock()

	if i != 1 {
		panic("fail to callback")
	}
}
