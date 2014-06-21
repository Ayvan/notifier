package models

import "time"

type Notice struct {
	Id       int
	Group    int
	Message  string

	Datetime time.Time
	Author   int
	//Group *Group
	//Author *User
}

func NewNotice(id int, group int, message string, datetime time.Time, author int) *Notice {

	return &Notice{id , group , message , datetime , author }
}
