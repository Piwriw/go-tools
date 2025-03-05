package compare

import "sort"

// Sortable 定义一个通用接口，支持排序 + 取元素 + 截取子序列
type Sortable[T any] interface {
	sort.Interface
	At(i int) T
	Slice(start, end int) Sortable[T] // 确保返回同样的类型
}

// 辅助函数：将 Sortable[T] 转换为 []T
func toSlice[T any, S Sortable[T]](s S) []T {
	result := make([]T, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.At(i)
	}
	return result
}

// TopNSort 泛型取前 N 名，支持并列
func TopNSort[T comparable, S Sortable[T]](items S, n int) []T {
	if items.Len() == 0 || n <= 0 {
		return nil
	}

	// 1. 降序排序
	sort.Sort(sort.Reverse(items))

	// 2. 取前 n 名的最后一个值
	if items.Len() <= n {
		return toSlice(items.Slice(0, items.Len()))
	}
	topN := toSlice(items.Slice(0, n))
	lastValue := items.At(n - 1)

	// 3. 继续遍历，找到所有相同的值
	for i := n; i < items.Len(); i++ {
		if items.At(i) == lastValue {
			topN = append(topN, items.At(i))
		} else {
			break
		}
	}
	return topN
}

// TopN 泛型取前N名，支持并列（前提已经排序好）
func TopN[T comparable](items []T, n int) []T {
	if len(items) == 0 {
		return nil
	}

	if len(items) <= n {
		return items
	}
	// 1. 取前 n 名的最后一个值
	topN := items[:n]
	lastValue := items[n-1]

	// 2. 继续遍历，找到所有相同的值
	for i := n; i < len(items); i++ {
		if items[i] == lastValue {
			topN = append(topN, items[i])
		} else {
			break
		}
	}
	return topN
}
