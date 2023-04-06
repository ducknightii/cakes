package main

import (
	"encoding/json"
	"fmt"
)

type Info struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func main() {
	m := make(map[string]string, 10)

	fmt.Printf("%+v %d\n", m, len(m))

	info := Info{
		ID:   1,
		Name: "aaaa",
	}

	body, err := json.Marshal(&info)
	fmt.Println(string(body), err)
}
