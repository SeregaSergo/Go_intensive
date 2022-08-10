package balance

import (
	t "github.com/SeregaSergo/Go_intensive/day_5/tree"
)

func recursiveCount(root *t.TreeNode) int {
	if root == nil {
		return 0
	}
	if root.HasToy {
		return 1 + recursiveCount(root.Right) + recursiveCount(root.Left)
	} else {
		return recursiveCount(root.Right) + recursiveCount(root.Left)
	}
}

func AreToysBalanced(root *t.TreeNode) bool {
	if root == nil {
		return true
	}
	if recursiveCount(root.Right) == recursiveCount(root.Left) {
		return true
	}
	return false
}
