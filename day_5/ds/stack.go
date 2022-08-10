package ds

import (
	t "github.com/SeregaSergo/Go_intensive/day_5/tree"
)

type Stack []*t.TreeNode

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Pop() (*t.TreeNode, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it by slicing it off.
		return element, true
	}
}

func (s *Stack) Push(elem *t.TreeNode) {
	if elem != nil {
		*s = append(*s, elem)
	}
}
