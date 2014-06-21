package models

type Notice struct {
	Id int
	Name string
	Message string
	Datetime string
	Group int
	Author int
	//Group *Group
	//Author *User
}
