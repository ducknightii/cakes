package main

import (
	"fmt"
	"time"
)

func main() {
	echo()
}

func echo() {
	fmt.Println("===")
	defer d()
	defer func() {
		d2()
	}()
	fmt.Println("===")
	time.Sleep(time.Second * 5)

	panic("aaa") // defer 依旧会执行
}

func d() {
	fmt.Println("ddd")
}

func d2() {
	fmt.Println("ddd2")
}
