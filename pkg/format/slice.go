package format

import (
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"unsafe"
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
	field, ok := elemType.FieldByName(orderBy)
	if !ok {
		return fmt.Errorf("field '%s' not found", orderBy)
	}

	fieldIndex := field.Index[0] // 直接获取索引，减少 FieldByName 调用

	// 构建 orderMap 提前索引
	orderMap := make(map[any]int, len(orderList))
	for i, v := range orderList {
		if !isComparable(v) {
			return fmt.Errorf("orderList contains an uncomparable value: %v", v)
		}
		orderMap[v] = i
	}

	sort.SliceStable(slice.Interface(), func(i, j int) bool {
		vi := slice.Index(i).Field(fieldIndex).Interface()
		vj := slice.Index(j).Field(fieldIndex).Interface()

		// 直接从 orderMap 取值，减少 map 查询次数
		iOrder, iExists := orderMap[vi]
		jOrder, jExists := orderMap[vj]

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

// SliceOrderByV2 对切片进行排序
// 会比v1稍微快一点
func SliceOrderByV2(rows any, orderBy string, orderList []any) error {
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
	field, ok := elemType.FieldByName(orderBy)
	if !ok {
		return fmt.Errorf("field '%s' not found", orderBy)
	}
	fieldOffset := field.Offset // 直接使用字段的偏移量

	// **提前解析 orderList 类型**
	orderMap := make(map[any]int, len(orderList))
	for i, v := range orderList {
		orderMap[v] = i
	}

	// 获取切片底层数据指针，避免反射带来的性能损耗
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(slice.UnsafeAddr()))
	basePtr := unsafe.Pointer(sliceHeader.Data)

	// 排序
	sort.SliceStable(slice.Interface(), func(i, j int) bool {
		// 计算第 i 和 j 个元素的内存地址
		ptrI := unsafe.Pointer(uintptr(basePtr) + uintptr(i)*uintptr(elemType.Size()))
		ptrJ := unsafe.Pointer(uintptr(basePtr) + uintptr(j)*uintptr(elemType.Size()))

		// 通过偏移量获取字段值
		fieldPtrI := unsafe.Pointer(uintptr(ptrI) + fieldOffset)
		fieldPtrJ := unsafe.Pointer(uintptr(ptrJ) + fieldOffset)

		// 获取字段的实际值（支持 string、int、float64）
		vi, vj := extractValue(field.Type, fieldPtrI, fieldPtrJ)

		// 直接从 orderMap 取值，减少 map 查询次数
		iOrder, iExists := orderMap[vi]
		jOrder, jExists := orderMap[vj]

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
		// 默认返回 false，保持原始顺序
		return false
	})

	return nil
}

// extractValue 提取字段值（支持 string、int、float64）
func extractValue(fieldType reflect.Type, fieldPtrI, fieldPtrJ unsafe.Pointer) (any, any) {
	switch fieldType.Kind() {
	case reflect.String:
		return *(*string)(fieldPtrI), *(*string)(fieldPtrJ)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return *(*int64)(fieldPtrI), *(*int64)(fieldPtrJ)
	case reflect.Float32, reflect.Float64:
		return *(*float64)(fieldPtrI), *(*float64)(fieldPtrJ)
	default:
		panic(fmt.Sprintf("unsupported field type: %v", fieldType))
	}
}
