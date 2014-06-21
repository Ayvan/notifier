package models
import "fmt"

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

func FindGroup(id string) *Group {
	return NewGroup(id, "grp", "1", []string{"1", "2", "3"})
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
