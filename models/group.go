package models

import (
	//"fmt"
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
	return &Group{id, name, owner, members}
}

func FindGroup(id string, redis services.Redis) *Group {
	//fmt.Println("FindGroup:", id)
	group := redis.Get(id)
	members := redis.GetMembers(id + ":members")

	var name string
	var author string

	for j := 0; j < len(group); j+=2 {
		switch group[j] {
		case "name":
			name = group[j+1]
		case "author":
			author = group[j+1]
		}
	}

	if len(group) >= 2 {
		return NewGroup(id, name, author, members)
	}

	return nil
}
