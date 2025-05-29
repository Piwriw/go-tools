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

	// 获取元素类型（处理指针情况）
	elemType := slice.Index(0).Type()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		elemType = elemType.Elem()
	}

	// 获取字段信息
	field, ok := elemType.FieldByName(orderBy)
	if !ok {
		return fmt.Errorf("field '%s' not found", orderBy)
	}
	fieldIndex := field.Index[0]

	// 构建 orderMap
	orderMap := make(map[any]int, len(orderList))
	for i, v := range orderList {
		if !isComparable(v) {
			return fmt.Errorf("orderList contains an uncomparable value: %v", v)
		}
		orderMap[v] = i
	}

	sort.SliceStable(slice.Interface(), func(i, j int) bool {
		// 获取元素值
		vi := slice.Index(i)
		vj := slice.Index(j)

		// 如果是指针，解引用
		if isPtr {
			vi = vi.Elem()
			vj = vj.Elem()
		}

		// 获取字段值
		viField := vi.Field(fieldIndex).Interface()
		vjField := vj.Field(fieldIndex).Interface()

		// 比较逻辑
		iOrder, iExists := orderMap[viField]
		jOrder, jExists := orderMap[vjField]

		if iExists && jExists {
			return iOrder < jOrder
		}
		if iExists {
			return true
		}
		if jExists {
			return false
		}
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

	// 获取切片元素类型（可能是结构体或结构体指针）
	elemType := slice.Index(0).Type()

	// 处理指针类型的情况
	var isPtr bool
	if elemType.Kind() == reflect.Ptr {
		isPtr = true
		elemType = elemType.Elem() // 获取指针指向的类型
	}

	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("slice elements must be struct or struct pointer")
	}

	field, ok := elemType.FieldByName(orderBy)
	if !ok {
		return fmt.Errorf("field '%s' not found", orderBy)
	}
	fieldOffset := field.Offset // 获取字段的内存偏移量

	// 提前解析 orderList 类型
	orderMap := make(map[any]int, len(orderList))
	for i, v := range orderList {
		orderMap[v] = i
	}

	// 获取切片数据指针
	dataPtr := unsafe.Pointer(slice.UnsafePointer())

	// 排序
	sort.SliceStable(slice.Interface(), func(i, j int) bool {
		// 计算第 i、j 个元素的地址
		ptrI := unsafe.Pointer(uintptr(dataPtr) + uintptr(i)*uintptr(slice.Index(0).Type().Size()))
		ptrJ := unsafe.Pointer(uintptr(dataPtr) + uintptr(j)*uintptr(slice.Index(0).Type().Size()))

		// 处理指针元素的情况
		if isPtr {
			ptrI = *(*unsafe.Pointer)(ptrI) // 解引用指针
			ptrJ = *(*unsafe.Pointer)(ptrJ) // 解引用指针
			if ptrI == nil || ptrJ == nil {
				// 处理nil指针的情况，这里将nil排在最后
				if ptrI == nil && ptrJ == nil {
					return false
				}
				return ptrJ == nil
			}
		}

		// 计算字段地址
		fieldPtrI := unsafe.Pointer(uintptr(ptrI) + fieldOffset)
		fieldPtrJ := unsafe.Pointer(uintptr(ptrJ) + fieldOffset)

		// 提取字段值
		vi, vj := extractValue(field.Type, fieldPtrI, fieldPtrJ)

		// 直接从 orderMap 取值
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
