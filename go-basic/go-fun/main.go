package main

import "fmt"

func main() {
	// slice 修改值
	s1 := []int{1, 2, 3, 4, 5}
	for _, item := range s1 {
		item++
		fmt.Printf("%p\n", &item) // 同一个地址空间
	}
	fmt.Println(s1) // [1 2 3 4 5]

	var v []*int
	for _, item := range s1 {
		v = append(v, &item)
	}
	for _, item := range v {
		fmt.Printf("%d\t", *item) // 拷贝过来的是同一地址空间 所以会是上一步 for range 的最后一个元素值 5	5	5	5	5
	}
	fmt.Println()

	var v2 []int
	for _, item := range s1 {
		v2 = append(v2, item)
	}
	s1[0]++ // 值拷贝过去v2 s1不会影响到 v2
	for _, item := range v2 {
		fmt.Printf("%d\t", item) // 1	2	3	4	5
	}
	fmt.Println()
}
