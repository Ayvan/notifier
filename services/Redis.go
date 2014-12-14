package services

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

type Redis struct {
	Host       string
	Port       string
	connection redis.Conn
}

func NewRedis(host string, port string) Redis {
	return Redis{host, port, nil}
}

func (this *Redis) Connect() {
	connection, error := redis.Dial("tcp", this.Host+":"+this.Port)

	if error != nil {
		log.Fatal(error)
	}

	this.connection = connection
}

func (this *Redis) Disconnect() {
	this.connection.Close()
}

func (this *Redis) AddNotice(noticeId string, noticeMessage string, noticeTime string) {
	this.connection.Send(
		"HMSET",
			"notice:"+noticeId,
		"noticeId",
		noticeId,
		"message",
		noticeMessage,
		"datetime",
		noticeTime,
	)

	this.connection.Flush()

	this.AddToRange("notices", "notice:"+noticeId, noticeTime)
}

func (this *Redis) Delete(key string) {
	this.connection.Send("DEL", key)
	this.connection.Flush()
}

func (this *Redis) AddToRange(rangeName string, key string, noticeTime string){
	this.connection.Send("ZADD", rangeName, noticeTime, key)

	this.connection.Flush()
}

func (this *Redis) DeleteFromRange(rangeName string, noticeId string) {
	this.connection.Send("ZREM", rangeName, noticeId)
	this.connection.Flush()
}

func (this *Redis) Get(key string) []string {

	result, error := this.connection.Do("HGETALL", key)

	if error != nil {
		fmt.Println(error)
		log.Fatal(error)
	}

	value, error := redis.Strings(result, error)

	if error != nil {
		fmt.Println(error)
		log.Fatal(error)
	}

	return value
}

func (this *Redis) GetRangeByScore(name string, min int, max int64) []string {

	result, error := this.connection.Do("ZRANGEBYSCORE", name, min, max)

	if error != nil {
		log.Fatal(error)
	}

	results, error := redis.Strings(result, error)

	return results
}

func (this *Redis) SearchKeys(query string) []string {

	result, error := this.connection.Do("KEYS", query)

	if error != nil {
		log.Fatal(error)
	}

	results, error := redis.Strings(result, error)

	if error != nil {
		log.Fatal(error)
	}

	return results
}
