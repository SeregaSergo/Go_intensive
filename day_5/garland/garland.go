package garland

import (
	"github.com/SeregaSergo/Go_intensive/day_5/ds"
	t "github.com/SeregaSergo/Go_intensive/day_5/tree"
)

func unrollLvl(s ds.Stack, dirToTheRight bool) []bool {
	coverdBranches := make([]bool, 0, len(s))
	nextLvlStack := make(ds.Stack, 0, 2)

	if s.IsEmpty() {
		return coverdBranches
	}

	for !s.IsEmpty() {
		elem, _ := s.Pop()
		coverdBranches = append(coverdBranches, elem.HasToy)
		if dirToTheRight {
			nextLvlStack.Push(elem.Left)
			nextLvlStack.Push(elem.Right)
		} else {
			nextLvlStack.Push(elem.Right)
			nextLvlStack.Push(elem.Left)
		}
	}

	return append(coverdBranches, unrollLvl(nextLvlStack, !dirToTheRight)...)
}

func UnrollGarland(root *t.TreeNode) []bool {
	s := make(ds.Stack, 0, 1)
	s = append(s, root)
	return unrollLvl(s, false)
}
