package main

import "fmt"

func main() {
	arr := []int{0, 1, 2, 3, 4}

	fmt.Println(arr)
	fmt.Println(arr[:0])
	arr = append(arr[:4], arr[5:]...)
	//fmt.Println(a1)
	//fmt.Println(arr)
	//a2 := append(arr[:4], arr[5:]...)
	//fmt.Println(a2)
	fmt.Println(arr, cap(arr), len(arr))
	fmt.Println(arr)

	str := "123"
	for _, c := range str {
		fmt.Printf("%T %d\n", c, c-'1')
	}

	fmt.Printf("%T: %v", str[1:], str[1:])
}
