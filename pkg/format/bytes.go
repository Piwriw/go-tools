package format

import "math"

// BytesToMB 将字节（Byte）转换为MB，保留两位小数
func BytesToMB(bytes float64) float64 {
	mb := bytes / (1024 * 1024)
	return math.Round(mb*100) / 100
}

// BytesToGB 将字节（Byte）转换为GB，保留两位小数
func BytesToGB(bytes float64) float64 {
	gb := bytes / (1024 * 1024 * 1024)
	return math.Round(gb*100) / 100
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
