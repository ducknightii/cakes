package main

import (
	"errors"
	"fmt"
)

func main() {
	var tErr = errors.New("aaa")
	t2 := fmt.Errorf("%w ttt", tErr)

	fmt.Println(errors.Unwrap(t2).Error())
	fmt.Println(t2.Error())
	t := errors.Unwrap(tErr)
	fmt.Println(t)

}
