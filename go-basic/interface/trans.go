package main

import "fmt"

func main() {
	var a []int = []int{1, 2, 3, 4}
	var t []interface{}
	for _, i := range a {
		t = append(t, i)
	}

	fmt.Printf("%+v", t)
}
