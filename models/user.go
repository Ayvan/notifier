package models

import "fmt"

type User struct {
	Id        string
	Name      string
}

func NewUser(id string, name string) *User {
	return &User{id , name}
}


func FindUser(id string) *User {
	name := fmt.Sprintf("user %d", id)
	return NewUser(id, name)
}
