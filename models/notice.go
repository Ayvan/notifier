package models

import (
	"notifier/services"
	"time"
	"strconv"
)

type Notice struct {
	Id      string
	Message string
	Datetime time.Time
}

func NewNotice(id string, message string, datetime time.Time) *Notice {

	return &Notice{id, message, datetime}
}

func NewNoticesFromRedis(redis services.Redis) []*Notice {

	currTime := time.Now().Unix()

	//Выбирает все записи, которые были до текущего времени
	results := redis.GetRangeByScore("notices", 0, currTime)
	notices := make([]*Notice, len(results), len(results))

	for i, noticeKey := range results {

		val := redis.Get(noticeKey)

		var noticeId string
		var message string
		var datetime string

		for j := 0; j < len(val); j+=2 {
			switch val[j] {
			case "noticeId":
				noticeId = val[j+1]
			case "message":
				message = val[j+1]
			case "datetime":
				datetime = val[j+1]
			}
		}

		intTime, err := strconv.ParseInt(datetime, 10, 64)
		var localTime time.Time
		if err == nil {
			localTime = time.Unix(intTime, 0)
		} else {
			localTime = time.Now()
		}

		notices[i] = NewNotice(noticeId, message, localTime)
	}

	return notices
}
