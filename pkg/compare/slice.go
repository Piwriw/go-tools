package compare

// Ordered 定义可排序的类型约束（用于 Min/Max）
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Comparable 定义可比较的类型约束（用于 Contains/Unique）
type Comparable interface {
	comparable
}

// Min 返回切片中的最小值
func Min[T Ordered](vals ...T) T {
	if len(vals) == 0 {
		var zero T
		return zero
	}
	minVal := vals[0]
	for _, v := range vals[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

// Max 返回切片中的最大值
func Max[T Ordered](vals ...T) T {
	if len(vals) == 0 {
		var zero T
		return zero
	}
	maxVal := vals[0]
	for _, v := range vals[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// Contains 判断切片中是否包含某个元素
func Contains[T comparable](vals []T, target T) bool {
	for _, v := range vals {
		if v == target {
			return true
		}
	}
	return false
}

// Unique 返回去重后的切片，保持原始顺序
func Unique[T comparable](vals []T) []T {
	if len(vals) == 0 {
		return make([]T, 0)
	}
	seen := make(map[T]struct{}, len(vals))
	result := make([]T, 0, len(vals))
	for _, v := range vals {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Chunk 将切片分成指定大小的块
func Chunk[T any](vals []T, size int) [][]T {
	if len(vals) == 0 || size <= 0 {
		return nil
	}
	if size >= len(vals) {
		return [][]T{vals}
	}

	chunks := make([][]T, 0, (len(vals)+size-1)/size)
	for i := 0; i < len(vals); i += size {
		end := i + size
		if end > len(vals) {
			end = len(vals)
		}
		chunks = append(chunks, vals[i:end])
	}
	return chunks
}
