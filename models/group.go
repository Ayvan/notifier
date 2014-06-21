package models


type Group struct {
	Id int
	Name string
	Owner int
	Members []int
	//Owner *User
	//Members []*User
}
