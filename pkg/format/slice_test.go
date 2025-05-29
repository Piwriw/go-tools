package format

import (
	"fmt"
	"testing"

	"github.com/mohae/deepcopy"

	"github.com/stretchr/testify/require"
)

type User struct {
	ID    int
	Name  string
	Role  string
	Score float64
}

func generateTestData(n int) []User {
	users := make([]User, n)
	roles := []string{"admin", "editor", "viewer", "guest"}
	for i := 0; i < n; i++ {
		users[i] = User{
			ID:    i + 1,
			Name:  fmt.Sprintf("User%d", i+1),
			Role:  roles[i%len(roles)],
			Score: float64(i%100) + 0.5,
		}
	}
	return users
}

func generateTestPointData(n int) []*User {
	users := make([]*User, n)
	roles := []string{"admin", "editor", "viewer", "guest"}
	for i := 0; i < n; i++ {
		users[i] = &User{
			ID:    i + 1,
			Name:  fmt.Sprintf("User%d", i+1),
			Role:  roles[i%len(roles)],
			Score: float64(i%100) + 0.5,
		}
	}
	return users
}

func TestBasicFunctionality2(t *testing.T) {
	users := generateTestPointData(5)
	orderBy := "Role"
	orderList := []any{"admin", "editor", "viewer", "guest"}
	for _, user := range users {
		fmt.Println(user)
	}
	fmt.Println("--------")
	// 测试V1
	if err := SliceOrderBy(&users, orderBy, orderList); err != nil {
		t.Errorf("SliceOrderByV1 failed: %v", err)
	}
	for _, user := range users {
		fmt.Println(user)
	}

}

func TestBasicFunctionality222(t *testing.T) {
	users := generateTestPointData(5)
	orderBy := "Role"
	orderList := []any{"admin", "editor", "viewer", "guest"}
	for _, user := range users {
		fmt.Println(user)
	}
	fmt.Println("--------")
	if err := SliceOrderByV2(&users, orderBy, orderList); err != nil {
		t.Errorf("SliceOrderByV1 failed: %v", err)
	}
	for _, user := range users {
		fmt.Println(user)
	}

}
func TestBasicFunctionality(t *testing.T) {
	users := generateTestData(100)
	orderBy := "Role"
	orderList := []any{"admin", "editor", "viewer", "guest"}

	// 测试V1
	err := SliceOrderBy(&users, orderBy, orderList)
	require.NoError(t, err)

	// 测试V2
	usersV2 := generateTestData(100)
	err = SliceOrderByV2(&usersV2, orderBy, orderList)
	require.NoError(t, err)

	// 验证结果是否一致
	for i := range users {
		require.Equal(t, users[i].Role, usersV2[i].Role)
	}
}
func BenchmarkSliceOrderBy(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000, 1000000}
	orderBy := "Role"
	orderList := []any{"admin", "editor", "viewer", "guest"}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			users := generateTestData(size)
			b.ResetTimer()

			b.Run("V1", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					usersCopy := make([]User, len(users))
					copy(usersCopy, users)
					_ = SliceOrderBy(&usersCopy, orderBy, orderList)
				}
			})

			b.Run("V2", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					usersCopy := make([]User, len(users))
					copy(usersCopy, users)
					_ = SliceOrderByV2(&usersCopy, orderBy, orderList)
				}
			})
		})
	}
}

func BenchmarkDifferentFieldTypes(b *testing.B) {
	size := 10000
	users := generateTestData(size)

	tests := []struct {
		name      string
		orderBy   string
		orderList []any
	}{
		{"StringField", "Role", []any{"admin", "editor", "viewer", "guest"}},
		{"IntField", "ID", []any{5, 3, 1, 2, 4}},
		{"FloatField", "Score", []any{50.5, 30.5, 10.5, 20.5, 40.5}},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.Run("V1", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					usersCopy := deepcopy.Copy(users).([]User)
					_ = SliceOrderBy(&usersCopy, tt.orderBy, tt.orderList)
				}
			})
			b.Run("V2", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					usersCopy := deepcopy.Copy(users).([]User)
					_ = SliceOrderByV2(&usersCopy, tt.orderBy, tt.orderList)
				}
			})
		})
	}
}
