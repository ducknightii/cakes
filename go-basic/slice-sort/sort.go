package main

import (
	"fmt"
	"sort"
)

type Item struct {
	val  int
	name string
}

func main() {
	arr := []Item{{
		val:  0,
		name: "a",
	}, {
		val:  2,
		name: "b",
	}, {
		val:  2,
		name: "c",
	}}

	sort.SliceStable(arr, func(i, j int) bool {
		return arr[i].val > arr[j].val
	})
	fmt.Println(arr)
}
