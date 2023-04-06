package main

import "fmt"

type User struct {
	ID   string
	Name string
}

func main() {
	u1 := User{
		ID:   "1",
		Name: "a",
	}

	u2 := User{
		ID:   "1",
		Name: "a",
	}

	fmt.Println(u1 == u2) // true
}
