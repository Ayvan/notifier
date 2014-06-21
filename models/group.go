package models
import "fmt"

type Group struct {
	Id      int
	Name    string
	Owner   int
	Members []int
	//Owner *User
	//Members []*User
}

func NewGroup(id int, name string, owner int, members []int) *Group {
	return &Group{id , name , owner , members }
}

func FindGroup(id int) *Group {
	return NewGroup(id, "grp", 1, []int{1, 2, 3})
}

func (this *Group) FindMembers() []*User {
	fmt.Println(this)
	ln := len(this.Members)
	users := make([]*User, ln, ln)
	for i, v := range this.Members {
		users[i] = FindUser(v)
	}
	return users

}
