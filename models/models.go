package models

//структуры в БД
type User struct {
	Id int
}

type Notice struct {
	Id int
}

type Group struct {
	Id int
}

//внутренние структуры для передачи между воркерами
type Message struct {
	Id int
}

type ChannelMessage struct {
	Id int
}

