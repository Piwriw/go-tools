package main

import "fmt"

func changeint(ints []int) {
	ints[0] = 100
}

func main() {
	ints := make([]int, 0)
	for i := 0; i < 10; i++ {
		ints = append(ints, i)
	}
	fmt.Println("Before", ints)
	changeint(ints)
	fmt.Println("After", ints)
}
