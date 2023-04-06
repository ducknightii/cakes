package main

import "fmt"

type Node struct {
	ID       string
	SubNodes []Node
}

func main() {
	n := t()
	fmt.Println(n == nil, len(n))

	nodes := []Node{
		{
			ID: "1",
			SubNodes: []Node{
				{
					ID: "1-1",
					SubNodes: []Node{
						{
							ID: "1-1-1",
						},
					},
				},
				{
					ID:       "1-2",
					SubNodes: []Node{},
				},
			},
		},
		{
			ID: "2",
			SubNodes: []Node{
				{
					ID: "2-1",
				},
			},
		},
		{
			ID: "3",
		},
	}

	fmt.Println(getTree(nodes, 3))
}

func t() []Node {
	n := new([]Node)
	return *n
}

func getTree(nodes []Node, level int) []Node {
	if level == 0 {
		return nodes
	}
	levelNode(&nodes, level)
	return nodes
}

type NodeListPtr *[]Node

func levelNode(nodeListPtr NodeListPtr, level int) {
	fmt.Printf("level:%d => %+v\n", level, nodeListPtr)
	if nodeListPtr == nil {
		return
	}
	if level == 0 {
		*nodeListPtr = nil
		return
	}
	for i := 0; i < len(*nodeListPtr); i++ {
		levelNode(&((*nodeListPtr)[i].SubNodes), level-1)
	}
}
