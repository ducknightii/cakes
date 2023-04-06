package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	str := ""
	fmt.Println(json.Valid([]byte(str)))
	str = "[]"
	fmt.Println(json.Valid([]byte(str)))
}
