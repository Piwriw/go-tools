package tree

import (
	"container/list"
	"fmt"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func NewTreeNode(v int) *TreeNode {
	return &TreeNode{
		Left:  nil, // 左子节点指针
		Right: nil, // 右子节点指针
		Val:   v,   // 节点值
	}
}
func (t *TreeNode) levelOrder(root TreeNode) []any {
	// 初始化队列，加入根节点
	queue := list.New()
	queue.PushBack(root)
	// 初始化一个切片，用于保存遍历序列
	nums := make([]any, 0)
	for queue.Len() > 0 {
		// 队列出队
		node := queue.Remove(queue.Front()).(*TreeNode)
		// 保存节点值
		nums = append(nums, node.Val)
		if node.Left != nil {
			queue.PushBack(node.Left)
		}
		if node.Right != nil {
			queue.PushBack(node.Right)
		}
	}
	return nums
}

/* 前序遍历 */
func preOrder(node *TreeNode) {
	if node == nil {
		return
	}
	// 访问优先级：根节点 -> 左子树 -> 右子树
	fmt.Printf("Current Node Value:%d\n", node.Val)
	preOrder(node.Left)
	preOrder(node.Right)
}

/* 中序遍历 */
func inOrder(node *TreeNode) {
	if node == nil {
		return
	}
	// 访问优先级：左子树 -> 根节点 -> 右子树
	inOrder(node.Left)
	fmt.Printf("Current Node Value:%d\n", node.Val)
	inOrder(node.Right)
}

/* 后序遍历 */
func postOrder(node *TreeNode) {
	if node == nil {
		return
	}
	// 访问优先级：左子树 -> 右子树 -> 根节点
	postOrder(node.Left)
	postOrder(node.Right)
	fmt.Printf("Current Node Value:%d\n", node.Val)
}
