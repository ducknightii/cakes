package main

import (
	"fmt"
	"sync"
	"time"
)

var poolIns sync.Pool

type ins string

func main() {
	poolIns.New = func() interface{} {
		fmt.Println("new....")
		return new(ins)
	}

	for i := 0; i < 20; i++ {
		for i := 0; i < 10; i++ {
			go func() {
				s := poolIns.Get()
				time.Sleep(time.Microsecond)
				poolIns.Put(s)
			}()
		}
		time.Sleep(time.Millisecond)
	}

	time.Sleep(time.Second)
}
