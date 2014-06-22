package models

import "time"
import (
	"iforgetgo/services"
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

	currTime :=time.Now().Unix()

	//Выбирает все записи, которые были до текущего времени
	results := redis.GetRangeByScore("notices", 0, currTime)
	notices := make([]*Notice, len(results), len(results))

	for i, noticeKey := range results {

		val := redis.Get(noticeKey)

		if len(val) >= 8 {

			group := val[1]
			message := val[3]
			//		time := val[5]
			author := val[7]

			notices[i] = NewNotice(noticeKey, group, message, time.Now(), author)
		}
	}

	return notices
}
