package main

import (
	"fmt"
	"math/big"
)

func main() {
	bigExample()

}
func bigExample() {
	// 创建大整数
	num1 := big.NewInt(123456789)
	num2 := big.NewInt(987654321)

	// 加法
	sum := new(big.Int).Add(num1, num2)
	fmt.Println("Sum:", sum)

	// 减法
	diff := new(big.Int).Sub(num1, num2)
	fmt.Println("Difference:", diff)

	// 乘法
	prod := new(big.Int).Mul(num1, num2)
	fmt.Println("Product:", prod)

	// 除法
	quot := new(big.Int).Div(num1, num2)
	fmt.Println("Quotient:", quot)

	// 取模
	rem := new(big.Int).Mod(num1, num2)
	fmt.Println("Remainder:", rem)
}
