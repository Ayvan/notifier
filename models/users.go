package models

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/go-av/curl"

)

type Users struct {
	NoticeId      string
	Users []User
}

func NewUsers(noticeId string, users []User) *Users {
	return &Users{noticeId, users}
}

func GetParticipants(noticeId string) *Users {

	// берем адрес из конфига
	url := beego.AppConfig.String("apiGetParticipants")+noticeId

	err, str := curl.String(url)

	if err != nil {
		return nil
	}

	//преобразуем строку JSON в массив байтов для декодера JSON
	jsonStr := []byte(str)

	var users []User

	// произведем разбор строки JSON в объявленный массив структур
	err = json.Unmarshal(jsonStr ,&users)

	if err != nil {
		return nil
	}

	if len(users) >= 1 {
		return NewUsers(noticeId, users)
	}

	return nil
}
