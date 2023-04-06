package main

import "fmt"

func main() {
	s := "cc"
	switch s {
	case "a":
		fmt.Println("a")
		return
	case "b":
		fmt.Println("b")
		return
	}

	fmt.Println("unknown")
}
