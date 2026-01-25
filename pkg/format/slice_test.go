package format

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// User 测试用结构体
type User struct {
	ID    int
	Name  string
	Role  string
	Score float64
}

// generateTestData 生成测试数据（值类型切片）
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

// generateTestPointData 生成测试数据（指针类型切片）
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

// TestSliceOrderBy 测试按字段自定义顺序排序切片
// 测试场景：正常场景、指针类型、空切片、不同字段类型
func TestSliceOrderBy(t *testing.T) {
	tests := []struct {
		name        string                        // 测试用例名称
		setupFunc   func() interface{}            // 生成测试数据
		orderBy     string                        // 排序字段
		orderList   []any                         // 自定义顺序
		wantErr     bool                          // 是否预期发生错误
		errContains string                        // 预期错误信息包含的字符串
		validate    func(*testing.T, interface{}) // 验证函数
	}{
		{
			name: "正常场景_指针切片按Role排序",
			setupFunc: func() interface{} {
				users := generateTestPointData(10)
				return &users
			},
			orderBy:   "Role",
			orderList: []any{"admin", "editor", "viewer", "guest"},
			wantErr:   false,
			validate: func(t *testing.T, data interface{}) {
				users := data.(*[]*User)
				// 验证顺序：admin -> editor -> viewer -> guest
				prevOrder := -1
				orderMap := map[string]int{"admin": 0, "editor": 1, "viewer": 2, "guest": 3}
				for _, u := range *users {
					currentOrder := orderMap[u.Role]
					assert.GreaterOrEqual(t, currentOrder, prevOrder, "顺序应非递减")
					prevOrder = currentOrder
				}
			},
		},
		{
			name: "正常场景_值切片按Role排序",
			setupFunc: func() interface{} {
				users := generateTestData(10)
				return &users
			},
			orderBy:   "Role",
			orderList: []any{"admin", "editor", "viewer", "guest"},
			wantErr:   false,
		},
		{
			name: "正常场景_按ID排序",
			setupFunc: func() interface{} {
				users := generateTestPointData(10)
				return &users
			},
			orderBy:   "ID",
			orderList: []any{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			wantErr:   false,
			validate: func(t *testing.T, data interface{}) {
				users := data.(*[]*User)
				// 验证ID顺序为降序
				for i := 0; i < len(*users)-1; i++ {
					assert.Greater(t, (*users)[i].ID, (*users)[i+1].ID)
				}
			},
		},
		{
			name: "边界条件_空切片",
			setupFunc: func() interface{} {
				users := []*User{}
				return &users
			},
			orderBy:   "Role",
			orderList: []any{"admin", "editor"},
			wantErr:   false,
		},
		{
			name: "边界条件_部分在orderList中",
			setupFunc: func() interface{} {
				users := generateTestPointData(10)
				return &users
			},
			orderBy:   "Role",
			orderList: []any{"admin", "viewer"}, // 只包含部分
			wantErr:   false,
			validate: func(t *testing.T, data interface{}) {
				users := data.(*[]*User)
				// admin和viewer应该在前面，editor和guest在后面
				foundOthers := false
				for _, u := range *users {
					if u.Role == "editor" || u.Role == "guest" {
						foundOthers = true
					}
					if u.Role == "admin" || u.Role == "viewer" {
						assert.False(t, foundOthers, "admin和viewer应该在editor和guest之前")
					}
				}
			},
		},
		{
			name: "异常场景_不是切片指针",
			setupFunc: func() interface{} {
				return "not a slice"
			},
			orderBy:     "Role",
			orderList:   []any{"admin"},
			wantErr:     true,
			errContains: "must be a pointer to a slice",
		},
		{
			name: "正常场景_单元素切片",
			setupFunc: func() interface{} {
				users := []*User{{ID: 1, Role: "admin"}}
				return &users
			},
			orderBy:   "Role",
			orderList: []any{"admin", "editor"},
			wantErr:   false,
		},
		{
			name: "边界条件_空orderList",
			setupFunc: func() interface{} {
				users := generateTestPointData(5)
				return &users
			},
			orderBy:   "Role",
			orderList: []any{},
			wantErr:   false,
			// 空orderList时，所有元素都不在列表中，应保持原顺序
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setupFunc()
			err := SliceOrderBy(data, tt.orderBy, tt.orderList)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains, "错误信息应包含指定字符串")
				}
			} else {
				require.NoError(t, err, "预期不发生错误")
				if tt.validate != nil {
					tt.validate(t, data)
				}
			}
		})
	}
}

// TestSliceOrderByV2 测试按字段自定义顺序排序切片V2版本
// 测试场景：正常场景、指针类型、空切片、不同字段类型、性能优化版本
func TestSliceOrderByV2(t *testing.T) {
	tests := []struct {
		name        string                        // 测试用例名称
		setupFunc   func() interface{}            // 生成测试数据
		orderBy     string                        // 排序字段
		orderList   []any                         // 自定义顺序
		wantErr     bool                          // 是否预期发生错误
		errContains string                        // 预期错误信息包含的字符串
		validate    func(*testing.T, interface{}) // 验证函数
	}{
		{
			name: "正常场景_指针切片按Role排序",
			setupFunc: func() interface{} {
				users := generateTestPointData(10)
				return &users
			},
			orderBy:   "Role",
			orderList: []any{"admin", "editor", "viewer", "guest"},
			wantErr:   false,
			validate: func(t *testing.T, data interface{}) {
				users := data.(*[]*User)
				// 验证顺序
				prevOrder := -1
				orderMap := map[string]int{"admin": 0, "editor": 1, "viewer": 2, "guest": 3}
				for _, u := range *users {
					currentOrder := orderMap[u.Role]
					assert.GreaterOrEqual(t, currentOrder, prevOrder, "顺序应非递减")
					prevOrder = currentOrder
				}
			},
		},
		{
			name: "正常场景_按ID排序",
			setupFunc: func() interface{} {
				users := generateTestPointData(10)
				return &users
			},
			orderBy:   "ID",
			orderList: []any{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			wantErr:   false,
		},
		{
			name: "边界条件_空切片",
			setupFunc: func() interface{} {
				users := []*User{}
				return &users
			},
			orderBy:   "Role",
			orderList: []any{"admin", "editor"},
			wantErr:   false,
		},
		{
			name: "异常场景_不是切片指针",
			setupFunc: func() interface{} {
				return "not a slice"
			},
			orderBy:     "Role",
			orderList:   []any{"admin"},
			wantErr:     true,
			errContains: "must be a pointer to a slice",
		},
		{
			name: "正常场景_值切片按Score排序",
			setupFunc: func() interface{} {
				users := generateTestData(10)
				return &users
			},
			orderBy:   "Score",
			orderList: []any{50.5, 30.5, 10.5, 20.5, 40.5},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setupFunc()
			err := SliceOrderByV2(data, tt.orderBy, tt.orderList)

			if tt.wantErr {
				require.Error(t, err, "预期发生错误")
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains, "错误信息应包含指定字符串")
				}
			} else {
				require.NoError(t, err, "预期不发生错误")
				if tt.validate != nil {
					tt.validate(t, data)
				}
			}
		})
	}
}

// TestSliceOrderByConsistency 测试V1和V2版本的一致性
func TestSliceOrderByConsistency(t *testing.T) {
	t.Run("V1和V2结果一致性", func(t *testing.T) {
		// 准备相同的测试数据
		data1 := generateTestData(100)
		data2 := generateTestData(100)

		orderBy := "Role"
		orderList := []any{"admin", "editor", "viewer", "guest"}

		// 使用V1排序
		err1 := SliceOrderBy(&data1, orderBy, orderList)
		require.NoError(t, err1)

		// 使用V2排序
		err2 := SliceOrderByV2(&data2, orderBy, orderList)
		require.NoError(t, err2)

		// 验证结果一致
		assert.Equal(t, len(data1), len(data2), "切片长度应一致")
		for i := range data1 {
			assert.Equal(t, data1[i].Role, data2[i].Role, "第%d个元素的Role应一致", i)
		}
	})

	t.Run("指针切片V1和V2结果一致性", func(t *testing.T) {
		data1 := generateTestPointData(100)
		data2 := generateTestPointData(100)

		orderBy := "Role"
		orderList := []any{"admin", "editor", "viewer", "guest"}

		err1 := SliceOrderBy(&data1, orderBy, orderList)
		require.NoError(t, err1)

		err2 := SliceOrderByV2(&data2, orderBy, orderList)
		require.NoError(t, err2)

		assert.Equal(t, len(data1), len(data2), "切片长度应一致")
		for i := range data1 {
			assert.Equal(t, data1[i].Role, data2[i].Role, "第%d个元素的Role应一致", i)
		}
	})
}

// TestSliceOrderByNilPointer 测试nil指针处理
func TestSliceOrderByNilPointer(t *testing.T) {
	t.Run("V2处理nil指针", func(t *testing.T) {
		users := []*User{
			{ID: 1, Role: "admin"},
			nil,
			{ID: 3, Role: "viewer"},
			nil,
			{ID: 5, Role: "guest"},
		}

		err := SliceOrderByV2(&users, "Role", []any{"admin", "viewer", "guest"})
		require.NoError(t, err)

		// 验证非nil元素按顺序排列
		var nonNilRoles []string
		for _, u := range users {
			if u != nil {
				nonNilRoles = append(nonNilRoles, u.Role)
			}
		}
		assert.Equal(t, []string{"admin", "viewer", "guest"}, nonNilRoles)
	})
}

// TestSliceOrderBy_Parallel 并发安全测试_SliceOrderBy
func TestSliceOrderBy_Parallel(t *testing.T) {
	t.Run("V1并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			users := generateTestPointData(10)
			err := SliceOrderBy(&users, "Role", []any{"admin", "editor", "viewer", "guest"})
			assert.NoError(t, err)
		}
	})

	t.Run("V2并发测试", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			users := generateTestPointData(10)
			err := SliceOrderByV2(&users, "Role", []any{"admin", "editor", "viewer", "guest"})
			assert.NoError(t, err)
		}
	})
}

// BenchmarkSliceOrderBy 性能基准测试_SliceOrderBy
func BenchmarkSliceOrderBy(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}
	orderBy := "Role"
	orderList := []any{"admin", "editor", "viewer", "guest"}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			users := generateTestData(size)
			b.ResetTimer()

			b.Run("V1", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					usersCopy := make([]User, len(users))
					copy(usersCopy, users)
					_ = SliceOrderBy(&usersCopy, orderBy, orderList)
				}
			})

			b.Run("V2", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					usersCopy := make([]User, len(users))
					copy(usersCopy, users)
					_ = SliceOrderByV2(&usersCopy, orderBy, orderList)
				}
			})
		})
	}
}

// BenchmarkDifferentFieldTypes 不同字段类型的性能基准测试
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
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					usersCopy := make([]User, len(users))
					copy(usersCopy, users)
					_ = SliceOrderBy(&usersCopy, tt.orderBy, tt.orderList)
				}
			})

			b.Run("V2", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					usersCopy := make([]User, len(users))
					copy(usersCopy, users)
					_ = SliceOrderByV2(&usersCopy, tt.orderBy, tt.orderList)
				}
			})
		})
	}
}

// BenchmarkPointerSlice 指针切片性能基准测试
func BenchmarkPointerSlice(b *testing.B) {
	size := 10000
	users := generateTestPointData(size)
	orderBy := "Role"
	orderList := []any{"admin", "editor", "viewer", "guest"}

	b.Run("V1_PointerSlice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			usersCopy := make([]*User, len(users))
			copy(usersCopy, users)
			_ = SliceOrderBy(&usersCopy, orderBy, orderList)
		}
	})

	b.Run("V2_PointerSlice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			usersCopy := make([]*User, len(users))
			copy(usersCopy, users)
			_ = SliceOrderByV2(&usersCopy, orderBy, orderList)
		}
	})
}

// ExampleSliceOrderBy 示例代码_SliceOrderBy
func ExampleSliceOrderBy() {
	// 定义用户结构体
	type User struct {
		ID   int
		Name string
		Role string
	}

	// 创建用户切片
	users := []User{
		{ID: 1, Name: "Alice", Role: "viewer"},
		{ID: 2, Name: "Bob", Role: "admin"},
		{ID: 3, Name: "Charlie", Role: "editor"},
		{ID: 4, Name: "David", Role: "viewer"},
	}

	// 按Role字段自定义顺序排序：admin -> editor -> viewer
	err := SliceOrderBy(&users, "Role", []any{"admin", "editor", "viewer"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 输出排序结果
	for _, u := range users {
		fmt.Printf("%s: %s\n", u.Name, u.Role)
	}

	// Output:
	// Bob: admin
	// Charlie: editor
	// Alice: viewer
	// David: viewer
}

// ExampleSliceOrderByV2 示例代码_SliceOrderByV2
func ExampleSliceOrderByV2() {
	// V2版本使用unsafe优化，性能更好
	type Product struct {
		ID       int
		Name     string
		Category string
	}

	products := []Product{
		{ID: 1, Name: "Laptop", Category: "Electronics"},
		{ID: 2, Name: "Book", Category: "Education"},
		{ID: 3, Name: "Phone", Category: "Electronics"},
		{ID: 4, Name: "Pen", Category: "Education"},
	}

	// 按Category自定义顺序排序
	err := SliceOrderByV2(&products, "Category", []any{"Education", "Electronics"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, p := range products {
		fmt.Printf("%s: %s\n", p.Name, p.Category)
	}

	// Output:
	// Book: Education
	// Pen: Education
	// Laptop: Electronics
	// Phone: Electronics
}
