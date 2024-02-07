package main

import (
	"fmt"
	"strings"
)

var a = []string{"a", "b", "c"}

/*
性能比较
strings.Join ≈ strings.Builder > bytes.Buffer > "+" > fmt.Sprintf
*/
func main() {
	addString()
	sprintfString()
	builderString()
	joinString()

}

// strings.join也是基于strings.builder来实现的,并且可以自定义分隔符，在join方法内调用了b.Grow(n)方法，这个是进行初步的容量分配，而前面计算的n的长度就是我们要拼接的slice的长度，因为我们传入切片长度固定，所以提前进行容量分配可以减少内存分配，很高效。
// strings.Join
func joinString() {
	join := strings.Join(a, "-")
	fmt.Println(join)
}

// builderString strings.Builder 拼接
func builderString() {
	var sb strings.Builder
	sb.WriteString(a[0])
	sb.WriteString(a[1])
	sb.WriteString(a[2])
	ret := sb.String()
	fmt.Println(ret)
}

// sprintfString Sprintf 字符串拼接
func sprintfString() {
	ret := fmt.Sprintf("%s%s%s", a[0], a[1], a[2])
	fmt.Println(ret)
}

// addString 字符串拼接
func addString() {

	//方式1：+
	ret := a[0] + a[1] + a[2]
	fmt.Println(ret)
}
