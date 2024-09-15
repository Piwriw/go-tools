package main

import (
	"fmt"
	"strconv"
)

type Person struct {
	Name    string
	Age     int
	Address []string
}

func changepeople(people []Person) {
	people[1].Address[0] = "Change"
}

func main() {
	peoples := make([]Person, 0)
	for i := 0; i < 10; i++ {
		peoples = append(peoples, Person{
			Name:    "joo" + strconv.Itoa(i),
			Age:     i,
			Address: []string{"\"address\",\"" + strconv.Itoa(i) + "\""},
		})
	}
	fmt.Println("Before", peoples[1], "Addr", peoples)
	peoplesCopy := make([]Person, len(peoples))
	copy(peoplesCopy, peoples)
	changepeople(peoplesCopy)
	fmt.Println("After", peoples[1], "Addr", peoples)
}
