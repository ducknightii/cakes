package main

import "fmt"

func main() {
	var a []*string
	var b []*string
	tmp := "aaa"
	a = append(a, &tmp)
	b = append(b, &tmp)
	fmt.Printf("a:%v b:%v\n", a, b)
	*a[0] = ""
	fmt.Printf("a:%v b:%v\n", a, b)

}
