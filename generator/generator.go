package main

import (
	"crypto/rand"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"os"
	"strconv"
	"time"
)

import m "math/rand"

func main() {

	firstNames := [50]string{"Marceline", "Shawnee", "Ned", "Yesenia", "Danika", "Randa", "Fernande", "Lenora", "Beatris", "Clifford", "Lynelle", "Shizuko", "Robbyn", "Genny", "Monroe", "Kattie", "Liz", "Dimple", "Merlin", "Vincenzo", "Joann", "Ciera", "Lakia", "Yon", "Drucilla", "Tandra", "Abby", "Lynne", "Edythe", "Debbie", "Karen", "Raven", "Merna", "Chi", "Tammy", "Altha", "Malika", "Nichole", "Jeannetta", "Joy", "Arletta", "Ying", "Blanch", "Jerlene", "Marvel", "Lizabeth", "Ambrose", "Tammie", "Trenton", "Renae"}
	connection, _ := redis.Dial("tcp", beego.AppConfig.String("redisHost")+":"+beego.AppConfig.String("redisPort"))

	fmt.Println("Generating data...")

	args := os.Args[1:]

	count, _ := strconv.Atoi(args[0])
	flush, _ := strconv.Atoi(args[1])

	// очищаем базу
	if flush == 1 {
		connection.Send("FLUSHDB")
	}

	for i := 0; i < count; i++ {

		// создаем NOTICE
		notice := randString(20)
		group := randString(15)
		user := randString(10)
		userFirstName := firstNames[m.Intn(len(firstNames))]
		timestamp := time.Now().Unix() + 10

		connection.Send(
			"HMSET",
			"notice:"+notice,
			"group",
			"group:"+group,
			"message",
			"this is message",
			"datetime",
			timestamp,
			"author",
			"user:"+user,
		)

		connection.Send("ZADD", "notices", timestamp, "notice:"+notice)

		// создаем GROUP
		connection.Send(
			"HMSET",
			"group:"+group,
			"author",
			"user:"+user,
			"name",
			"group "+strconv.Itoa(i+1),
		)

		// добавляем пользователя в созданную группу
		connection.Send(
			"SADD",
			"group:"+group+":members",
			"user:"+user,
		)

		// создаем пользователя
		connection.Send(
			"HMSET",
			"user:"+user,
			"name",
			userFirstName,
		)

		// добавляем адреса пользователю
		connection.Send(
			"HMSET",
			"user:"+user+":addresses",
			"email",
			"test@example.com",
			"phone",
			"00000000",
		)

	}

	connection.Flush()

	fmt.Println("Generated records:", count)

}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
