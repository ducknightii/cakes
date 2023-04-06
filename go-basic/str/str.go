package main

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

func main() {
	s := "嗷嗷aaa"

	// 是以 rune 进行分割的
	for _, c := range s {
		fmt.Printf("%c %T\n", c, c)
	}

	var g errgroup.Group
	for i := 0; i < 10; i++ {
		seconds := time.Second * time.Duration(i)
		g.Go(func() error {
			time.Sleep(seconds)
			if seconds < time.Second*2 {
				return errors.New("too short")
			}

			return nil
		})
	}

	now := time.Now()
	err := g.Wait()
	fmt.Printf("%ds: %v\n", time.Now().Sub(now)/time.Second, err)

	ch := make(chan int)
}
