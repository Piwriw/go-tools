package stack

type arrayStack struct {
	data []int
}

/* 初始化栈 */
func newArrayStack() *arrayStack {
	return &arrayStack{
		// 设置栈的长度为 0，容量为 16
		data: make([]int, 0, 16),
	}
}

/* 栈的长度 */
func (s *arrayStack) size() int {
	return len(s.data)
}

/* 栈是否为空 */
func (s *arrayStack) isEmpty() bool {
	return s.size() == 0
}

/* 入栈 */
func (s *arrayStack) push(value int) {
	s.data = append(s.data, value)
}

/* 出栈 */
func (s *arrayStack) pop() any {
	val := s.peek()
	s.data = s.data[:len(s.data)-1]
	return val
}
func (s *arrayStack) peek() any {
	if s.isEmpty() {
		return nil
	}
	val := s.data[len(s.data)-1]
	return val
}

/* 获取 toSlice 用于打印 */
func (s *arrayStack) toSlice() []int {
	return s.data
}
