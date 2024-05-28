package tree

import "testing"

func Test(t *testing.T) {
	root := NewTreeNode(0)
	left := NewTreeNode(1)
	right := NewTreeNode(2)
	left2 := NewTreeNode(3)
	root.Left = left
	root.Left.Left = left2
	root.Right = right
	t.Log("preOrder:")
	preOrder(root)
	t.Log("inOrder:")
	inOrder(root)
	t.Log("postOrder:")
	postOrder(root)
}
