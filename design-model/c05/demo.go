package main

import (
	"fmt"
	resource_pool "github.com/ducknightii/cakes/design-model/c05/builder/resource-pool"
)

func main() {
	b := resource_pool.Builder{}
	b.SetName("aaa")

	p, _ := b.Build()

	fmt.Printf("%T %+v\n", p, p)
}
