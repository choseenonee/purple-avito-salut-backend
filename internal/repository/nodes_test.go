package repository

import (
	"fmt"
	"testing"
	"time"
)

func TestNodesRepo_GetMicrocategoryTree(t *testing.T) {
	db := initDB()

	nodesRepo := InitNodesRepo(db)

	_, _ = nodesRepo.GetMicrocategoryTree()
}

func recursive(index int, in [][4]int, ans []int, lastWithPrice int, lastIndex int, isFirst bool) bool {
	isFound := false
	for _, i := range in[lastIndex:] {
		if isFirst && i[0] != index {
			break
		}
		if i[0] != index {
			continue
		}
		if i[2] != 0 {
			lastWithPrice = index
		}
		if !isFound {
			ans[index-1] = lastWithPrice
		}
		if !recursive(i[1], in, ans, lastWithPrice, index-1, false) {
			if i[3] != 0 {
				ans[i[1]-1] = i[1]
			} else {
				ans[i[1]-1] = lastWithPrice
			}
		}
		isFound = true
		//_ = g
	}
	return isFound
}

func TestNodesRepo_Algo(t *testing.T) {
	var mean int

	for i := 0; i < 1000; i++ {
		a := time.Now()

		in := [][4]int{{1, 2, 100, 0}, {1, 3, 100, 300}, {1, 4, 100, 400}, {2, 5, 0, 0}, {2, 6, 0, 0}, {3, 7, 300, 0}, {5, 8, 0, 900}, {5, 9, 0, 0}}

		ans := make([]int, 9)

		recursive(1, in, ans, 1, 0, true)
		for index, value := range ans {
			fmt.Println(fmt.Sprintf("from %v to %v", index+1, value))
		}

		fmt.Println(time.Since(a))

		fmt.Println(ans)
	}
}

var tree = map[int][]int{}
var bools = map[int]bool{}

type Row struct {
	Parent int
	Child  int
	CPrice bool
}

func treeRunner(start []int, prevPriceID int) {
	var nextID int
	for _, child := range start {
		if bools[child] {
			fmt.Println(child, child)
			nextID = child
		} else {
			fmt.Println(child, prevPriceID)
			nextID = prevPriceID
		}

		treeRunner(tree[child], nextID)
	}
}

func TestNodesRepo_AlgoKolya(t *testing.T) {
	a := time.Now()
	data := []Row{
		{
			1,
			2,
			false,
		},
		{
			1,
			3,
			true,
		},
		{
			1,
			4,
			true,
		},
		{
			2,
			5,
			false,
		},
		{
			2,
			6,
			false,
		},
		{
			5,
			8,
			true,
		},
		{
			5,
			9,
			false,
		},
		{
			3,
			7,
			false,
		},
	}

	for _, row := range data {
		tree[row.Parent] = append(tree[row.Parent], row.Child)
		bools[row.Child] = row.CPrice
	}

	fmt.Println(1, 1)

	treeRunner(tree[1], 1)
	fmt.Println(time.Since(a))
}
