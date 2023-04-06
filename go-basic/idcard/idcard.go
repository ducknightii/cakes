package main

import (
	"fmt"
	"time"
)

func main() {
	idCard := "123122199621022811"
	dateStr := idCard[6:14]
	fmt.Println(dateStr)

	date, _ := time.ParseInLocation("20060102", dateStr, time.Local)
	fmt.Println(date)
}
