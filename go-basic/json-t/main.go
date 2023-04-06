package main

import (
	"encoding/json"
	"fmt"
)

type T struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func main() {
	t := T{
		A: "aaa",
		B: 111,
	}

	dataBytes, _ := json.Marshal(t)

	var res interface{}
	err := json.Unmarshal(dataBytes, &res)

	fmt.Printf("res:%+v, err:%v\n", res, err)

	dataBytes, err = json.Marshal(res)
	fmt.Printf("res:%s, err:%v\n", dataBytes, err)

}
