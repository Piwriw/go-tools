package utils

import (
	"fmt"
	"github.com/huandu/go-clone"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Person struct {
	Name    string
	Age     int
	Address []string
}

var wg sync.WaitGroup

func fixpeople(people []Person) {
	people[1].Address[0] = "Change"
}
func TestCopy(t *testing.T) {
	peoples := make([]Person, 0)
	for i := 0; i < 1000000; i++ {
		peoples = append(peoples, Person{
			Name:    "joo" + strconv.Itoa(i),
			Age:     i,
			Address: []string{"address" + strconv.Itoa(i), strconv.Itoa(i)},
		})
	}
	wg.Add(3)
	go func([]Person) {
		startDeepCOpy(peoples)
		wg.Done()
	}(peoples)
	go func([]Person) {
		startClone(peoples)
		wg.Done()
	}(peoples)
	go func() {
		startCloneSlow(peoples)
		wg.Done()
	}()

	//t.Log(peoples)
	wg.Wait()
}

func startCloneSlow(people []Person) []Person {
	now := time.Now()
	slowly := clone.Slowly(people).([]Person)
	fixpeople(slowly)

	fmt.Println("startCloneSlow Time", time.Since(now).Seconds())
	return slowly
}
func startClone(people []Person) []Person {
	now := time.Now()
	peoples := clone.Clone(people).([]Person)
	fixpeople(peoples)

	fmt.Println("startClone Time", time.Since(now).Seconds())
	return people
}

func startDeepCOpy(peoplesCopy []Person) []Person {
	now := time.Now()

	people := Copy(peoplesCopy).([]Person)
	fixpeople(people)
	fmt.Println("startDeepCopy Time", time.Since(now).Seconds())
	return people
}
