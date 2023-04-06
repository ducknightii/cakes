package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println(now.String())

	str := "2022-07-09T10:38:18.655137347+08:00"
	p, _ := time.ParseInLocation(time.RFC3339, str, time.Local)
	fmt.Println(p)

	min := 121
	fmt.Printf("%02d:%02d", min/60, min%60)
}
