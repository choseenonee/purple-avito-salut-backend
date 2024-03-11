package models

type NodeRaw struct {
	ID          int
	HasPrice    bool
	ChildrenIDs []int
	Parent      *Node
}

type Node struct {
	ID       int
	HasPrice bool
	Children []*Node
	Parent   *Node
}
