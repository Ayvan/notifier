package models

import (
	"iforgetgo/services"
	"time"
	"strconv"
)

type Notice struct {
	Id      string
	Group   string
	Message string

	Datetime time.Time
	Author   string
	//Group *Group
	//Author *User
}

func NewNotice(id string, group string, message string, datetime time.Time, author string) *Notice {

	return &Notice{id, group, message, datetime, author}
}

func NewNoticesFromRedis(redis services.Redis) []*Notice {

	currTime := time.Now().Unix()

	//Выбирает все записи, которые были до текущего времени
	results := redis.GetRangeByScore("notices", 0, currTime)
	notices := make([]*Notice, len(results), len(results))

	for i, noticeKey := range results {

		val := redis.Get(noticeKey)

		var group string
		var message string
		var datetime string
		var author string

		for j := 0; j < len(val); j+=2 {
			switch val[j] {
			case "group":
				group = val[j+1]
			case "message":
				message = val[j+1]
			case "datetime":
				datetime = val[j+1]
			case "author":
				author = val[j+1]
			}
		}

		intTime, err := strconv.ParseInt(datetime, 10, 64)
		var localTime time.Time
		if err == nil {
			localTime = time.Unix(intTime, 0)
		} else {
			localTime = time.Now()
		}

		notices[i] = NewNotice(noticeKey, group, message, localTime, author)
	}

	return notices
}
