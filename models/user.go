package models

import "fmt"

type User struct {
	Id    int
	Name  string
	Phone string
	Mail  string
}

func NewUser(id int, name string, phone string, mail string) *User {
	return &User{id , name , phone , mail }
}


func FindUser(id int) *User {
	name := fmt.Sprintf("user %d", id)
	return NewUser(id, name, "79261112233", name+"@iforget.biz")
}
