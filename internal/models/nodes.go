package models

type Node struct {
	ID       int
	HasPrice bool
	Children []*Node
	Parent   *Node
}
