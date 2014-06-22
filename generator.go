package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"crypto/rand"
	"time"
)

func main() {

	connection, _ := redis.Dial("tcp", beego.AppConfig.String("redisHost")+":"+beego.AppConfig.String("redisPort"))

	fmt.Println("Generating data...")

	args := os.Args[1:]

	count, _ := strconv.Atoi(args[0])
	flush, _ := strconv.Atoi(args[1])

	// очищаем базу
	if (flush == 1) {
		connection.Send("FLUSHDB")
	}

	for i := 0 ; i < count; i++ {

		// создаем NOTICE
		notice := randString(20)
		group := randString(15)
		user := randString(10)
		timestamp := time.Now().Unix() + 10

		connection.Send(
			"HMSET",
				"notice:"+notice,
			"group",
				"group:"+group,
			"message",
			"this is message",
			"time",
			timestamp,
			"author",
				"user:"+user,
		)

		connection.Send("ZADD", "notices", timestamp, "notice:"+notice)

		// создаем GROUP
		connection.Send(
			"HMSET",
				"groups:"+group,
			"owner",
				"user:"+user,
			"name",
				"group "+strconv.Itoa(i+1),
		)

		// добавляем пользователя в созданную группу
		connection.Send(
			"SADD",
					"groups:"+group+":members",
				"user:"+user,
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
