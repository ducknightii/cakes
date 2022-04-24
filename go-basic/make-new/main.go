package main

func main() {
	// make 内置函数 只能用来创建 slice map chan，返回的变量类型为传入Type 而不是指针
	s := make([]int, 0, 10)
}
