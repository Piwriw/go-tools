package format

import "math"

// BytesToMBExact 将字节（Byte）转换为MB，保留小数
func BytesToMBExact[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](bytes T) float64 {
	return float64(bytes) / (1024 * 1024)
}

// BytesToGBExact 将字节（Byte）转换为GB，保留小数
func BytesToGBExact[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](bytes T) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}

// BytesToMB 将字节（Byte）转换为MB，保留两位小数
func BytesToMB[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](bytes T) T {
	mb := float64(bytes) / (1024 * 1024)
	return T(math.Round(mb*100) / 100)
}

// BytesToGB 将字节（Byte）转换为GB，保留两位小数
func BytesToGB[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](bytes T) T {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return T(math.Round(gb*100) / 100)
}

// KBToMB 将KB转换为MB，保留两位小数
func KBToMB(kb float64) float64 {
	mb := kb / 1024
	return math.Round(mb*100) / 100
}

// KBToGB 将KB转换为GB，保留两位小数
func KBToGB(kb float64) float64 {
	gb := kb / (1024 * 1024)
	return math.Round(gb*100) / 100
}
