package models

import "fmt"

type User struct {
	Id    string
	Name  string
	Phone string
	Mail  string
}

func NewUser(id string, name string, phone string, mail string) *User {
	return &User{id , name , phone , mail }
}


func FindUser(id string) *User {
	name := fmt.Sprintf("user %d", id)
	return NewUser(id, name, "79261112233", name+"@iforget.biz")
}
