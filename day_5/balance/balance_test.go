package balance

import (
	tr "github.com/SeregaSergo/Go_intensive/day_5/tree"
	"testing"
)

func TestAreToysBalanced(t *testing.T) {
	tree1 := &tr.TreeNode{
		false,
		&tr.TreeNode{
			false,
			&tr.TreeNode{false, nil, nil},
			&tr.TreeNode{true, nil, nil},
		},
		&tr.TreeNode{
			true,
			nil,
			nil,
		},
	}
	tree2 := &tr.TreeNode{
		true,
		&tr.TreeNode{
			true,
			&tr.TreeNode{true, nil, nil},
			&tr.TreeNode{false, nil, nil},
		},
		&tr.TreeNode{
			false,
			&tr.TreeNode{true, nil, nil},
			&tr.TreeNode{true, nil, nil},
		},
	}
	tree3 := &tr.TreeNode{
		true,
		&tr.TreeNode{true, nil, nil},
		&tr.TreeNode{false, nil, nil},
	}
	tree4 := &tr.TreeNode{
		false,
		&tr.TreeNode{
			true,
			nil,
			&tr.TreeNode{true, nil, nil},
		},
		&tr.TreeNode{
			false,
			nil,
			&tr.TreeNode{true, nil, nil},
		},
	}
	if AreToysBalanced(tree1) != true {
		t.Error("Expected true in tree1")
	}
	if AreToysBalanced(tree2) != true {
		t.Error("Expected true in tree2")
	}
	if AreToysBalanced(tree3) != false {
		t.Error("Expected false in tree3")
	}
	if AreToysBalanced(tree4) != false {
		t.Error("Expected false in tree4")
	}
}
