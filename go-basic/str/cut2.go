package main

import (
	"fmt"
	"strings"
)

// "0-0-0"
func main() {
	strArr := []string{
		"0-0-0",
		"0-0",
		"1",
		"2-1-2",
		"3-1",
	}

	var idMap = make(map[string][]string)
	for _, item := range strArr {
		if _, ok := idMap[item]; !ok {
			idMap[item] = []string{}
		}
		cuts := strings.Split(item, "-")
		for i := len(cuts) - 1; i > 0; i-- {
			key := strings.Join(cuts[:i], "-")
			value := strings.Join(cuts[:i+1], "-")
			idMap[key] = append(idMap[key], value)
		}
	}

	for key, val := range idMap {

		fmt.Printf("%s => %+v\n", key, val)

	}

	str := "http://asdasdasdas/qr/asdas+="
	cuts := strings.Split(str, "/qr/")
	fmt.Printf("cuts:%+v\n", cuts)
}
