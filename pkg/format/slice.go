package format

import (
	"fmt"
	"log/slog"
	"reflect"
	"sort"
)

func isComparable(value any) bool {
	return reflect.TypeOf(value).Comparable()
}

// SliceOrderBy 根据字段名和自定义顺序排序
// 需要传入指针切片，例如：[]*User
// 不匹配就按照原来的顺序排序
func SliceOrderBy(rows any, orderBy string, orderList []any) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("SliceOrderBy Panic:", slog.Any("err", r))
		}
	}()
	value := reflect.ValueOf(rows)

	// 确保传入的是 *slice
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("rows must be a pointer to a slice")
	}

	slice := value.Elem()
	if slice.Len() == 0 {
		return nil
	}

	elemType := slice.Index(0).Type()
	_, ok := elemType.FieldByName(orderBy)
	if !ok {
		return fmt.Errorf("field '%s' not found", orderBy)
	}

	// 记录 orderList 的索引，方便查找
	orderMap := make(map[any]int)

	for i, v := range orderList {
		if !isComparable(v) {
			return fmt.Errorf("orderList contains an uncomparable value: %v", v)
		}
		orderMap[v] = i
	}

	// 解析字段类型
	sort.Slice(slice.Interface(), func(i, j int) bool {
		vi := slice.Index(i).FieldByName(orderBy)
		vj := slice.Index(j).FieldByName(orderBy)

		// 处理 orderList 优先排序
		viKey := vi.Interface()
		vjKey := vj.Interface()

		iOrder, iExists := orderMap[fmt.Sprintf("%v", viKey)]
		jOrder, jExists := orderMap[fmt.Sprintf("%v", vjKey)]

		// 如果都在 orderList 里，按顺序排序
		if iExists && jExists {
			return iOrder < jOrder
		}

		// 仅 i 在 orderList 里，i 优先
		if iExists {
			return true
		}

		// 仅 j 在 orderList 里，j 优先
		if jExists {
			return false
		}
		// 不匹配，就按照原来的顺序排序
		return false
	})

	return nil
}
