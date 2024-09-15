package main

import (
	"fmt"
	"strconv"
)

type Person struct {
	Name    string
	Age     int
	Address string
}

func changepeople(people [10]*Person) {
	people[1].Address = "Change"
}

func main() {
	var peoples [10]*Person
	for i := 0; i < 10; i++ {
		peoples[i] = &Person{
			Name:    "joo" + strconv.Itoa(i),
			Age:     i,
			Address: "address",
		}
	}
	fmt.Println("Before", peoples[1], "Addr", peoples)
	changepeople(peoples)
	fmt.Println("After", peoples[1], "Addr", peoples)
}
