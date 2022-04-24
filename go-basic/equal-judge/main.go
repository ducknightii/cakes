package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	//array()
	//slice()
	interf()
}

func array() {
	// 数组
	a := [2]interface{}{1, 2}
	a[1] = a       // 值拷贝
	fmt.Println(a) // [1 [1 2]]
}

func slice() {
	// slice
	sl := []interface{}{1, 2}
	sl[1] = sl                                         // 浅拷贝（地址），会造成循环引用 !!!
	fmt.Printf("sl cap:%d len:%d\n", cap(sl), len(sl)) //  cap:2 len:2
	fmt.Println(sl)                                    // panic: fatal error: stack overflow
}

func interf() {
	var f *os.File
	var r io.Reader = f
	var rc io.ReadCloser = f
	// 接口值的比较不要求接口类型（注意不是动态类型）完全相同，只要一个接口可以转化为另一个就可以比较。
	fmt.Println(r == rc) // true
	fmt.Println(rc == r) // true
}
