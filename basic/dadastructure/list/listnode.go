package list

type ListNode struct {
	Val  int
	Next *ListNode
}

func NewListNode(val int) *ListNode {
	return &ListNode{
		Val:  val,
		Next: nil,
	}
}

/* 在链表的节点 n0 之后插入节点 P */
func insertNode(n0 *ListNode, P *ListNode) {
	n1 := n0.Next
	P.Next = n1
	n0.Next = P
}

// 删除链表中特定值的节点
func deleteNodeWithValue(head *ListNode, val int) *ListNode {
	// 如果要删除的节点是头节点
	if head.Val == val {
		return head.Next
	}
	current := head
	for head.Next != nil {
		if current.Next.Val == val {
			current.Next = current.Next.Next
			break
		}
		current = current.Next
	}
	return head
}

/* 访问链表中索引为 index 的节点 */
func access(head *ListNode, index int) *ListNode {
	for i := 0; i < index; i++ {
		if head == nil {
			return nil
		}
		head = head.Next
	}
	return head
}

/* 查找某索引的节点 */
func searchNode(head *ListNode, target int) int {
	index := 0
	for head != nil {
		if head.Val == target {
			return index
		}
		head = head.Next
		index++
	}
	return -1
}
