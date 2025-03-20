package format

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func Test(t *testing.T) {
	users := []any{
		[]byte("A"),
	}
	// 自定义排序顺序
	customOrder := []any{"C", "B"}
	err := SliceOrderBy(&users, "Name", customOrder)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(users)
	}
}

func TestOrder(t *testing.T) {
	users := []User{
		{"B", 22},
		{"C", 30},
		{"D", 30},
		{"A", 25},

		{"E", 55},
	}
	// 自定义排序顺序
	customOrder := []any{"C", "B"}
	err := SliceOrderBy(&users, "Name", customOrder)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(users)
	}
}
