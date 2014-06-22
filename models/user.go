package models

import (
	"iforgetgo/services"
)

type User struct {
	Id        string
	Name      string
}

func NewUser(id string, name string) *User {
	return &User{id , name}
}


func FindUser(id string, redis services.Redis) *User {

	result := redis.Get(id)
	if len(result) >= 2 {
		userName := result[1]
		return NewUser(id, userName)
	}
	return nil
}
