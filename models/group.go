package models

import (
	"fmt"
	"iforgetgo/services"
)

type Group struct {
	Id      string
	Name    string
	Owner   string
	Members []string
	//Owner *User
	//Members []*User
}

func NewGroup(id string, name string, owner string, members []string) *Group {
	return &Group{id , name , owner , members }
}

func FindGroup(id string, redis *services.Redis) *Group {
	fmt.Println("FindGroup:", id)
	val := redis.Get(id)
	members := redis.GetMembers(id + ":members")
	return NewGroup(val[0], val[1], val[2], members)
}

