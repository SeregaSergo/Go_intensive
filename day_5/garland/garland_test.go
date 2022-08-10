package garland

import (
	tr "github.com/SeregaSergo/Go_intensive/day_5/tree"
	"reflect"
	"testing"
)

func TestUnrollGarland(t *testing.T) {
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
	if !reflect.DeepEqual(UnrollGarland(tree1), []bool{false, false, true, true, false}) {
		t.Error("Tree1:", UnrollGarland(tree1))
	}
	if !reflect.DeepEqual(UnrollGarland(tree2), []bool{true, true, false, true, true, false, true}) {
		t.Error("Tree2: ", UnrollGarland(tree2))
	}
	if !reflect.DeepEqual(UnrollGarland(tree3), []bool{true, true, false}) {
		t.Error("Tree3: ", UnrollGarland(tree3))
	}
	if !reflect.DeepEqual(UnrollGarland(tree4), []bool{false, true, false, true, true}) {
		t.Error("Tree4: ", UnrollGarland(tree4))
	}
}
