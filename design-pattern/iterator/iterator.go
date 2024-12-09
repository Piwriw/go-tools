package main

import "fmt"

type User struct {
	name string
	age  int
}

func (u *User) String() string {
	return fmt.Sprintf("Name:%s Age:%d", u.name, u.age)
}

type iterator interface {
	hasNext() bool
	next() *User
}

type userIterator struct {
	users []*User
	index int
}

func (u *userIterator) hasNext() bool {
	return u.index < len(u.users)
}

func (u *userIterator) next() *User {
	if u.hasNext() {
		user := u.users[u.index]
		u.index++
		return user
	}
	return nil
}

type userCollection struct {
	users []*User
}

func (u *userCollection) createIterator() iterator {
	return &userIterator{
		users: u.users,
	}
}

func main() {
	userK := &User{
		name: "Kevin",
		age:  30,
	}
	userD := &User{
		name: "Diamond",
		age:  25,
	}

	userCollection := &userCollection{
		users: []*User{userK, userD},
	}
	iterator := userCollection.createIterator()
	for iterator.hasNext() {
		user := iterator.next()
		fmt.Println(user)
	}
}
