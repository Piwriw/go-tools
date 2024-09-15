package utils

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/huandu/go-clone"
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
	for i := 0; i < 10; i++ {
		peoples = append(peoples, Person{
			Name:    "joo" + strconv.Itoa(i),
			Age:     i,
			Address: []string{"address" + strconv.Itoa(i), strconv.Itoa(i)},
		})
	}
	clone.Slowly(peoples)

	//peoplesCopy := make([]Person, len(peoples))
	//copy(peoplesCopy, peoples)
	//peoplesDeepCopy := Copy(peoplesCopy).([]Person)

	wg.Add(3)
	go func([]Person) {
		fmt.Println(startDeepCOpy(peoples))
		wg.Done()
	}(peoples)
	go func([]Person) {
		fmt.Println(startClone(peoples))
		wg.Done()
	}(peoples)
	go func() {
		fmt.Println(startCloneSlow(peoples))
		wg.Done()
	}()

	t.Log(peoples)
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
	fmt.Println("startDeepCOpy Time", time.Since(now).Seconds())
	return people
}
