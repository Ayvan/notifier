package services

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type Redis struct {
	Host     string
	Port     string

	connection redis.Conn
}

func (this *Redis) Connect() {
	connection, error := redis.Dial("tcp", this.Host+":"+this.Port);

	if (error != nil) {
		log.Fatal(error)
	}

	this.connection = connection
}

func (this *Redis) Delete(key string) {
	this.connection.Send("DEL", key);
}
