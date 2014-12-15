package models

import (
	"github.com/astaxie/beego"
	"notifier/services/curl"
	"encoding/json"
	"crypto/sha1"
	"crypto/hmac"
	"fmt"
	"encoding/base64"
)

type User struct {
	Name string
	Addresses []UserAddresses
}


type UserAddresses struct {
	Channel string
	Address string
}


type Users struct {
	NoticeId      string
	Users []User
}

func NewUsers(noticeId string, users []User) *Users {
	return &Users{noticeId, users}
}

func GetParticipants(noticeId string) *Users {

	key := []byte(beego.AppConfig.String("apiSecretKey"))
	mac := hmac.New(sha1.New,key)
	macMessage := []byte("GetParticipants"+noticeId)
	mac.Write(macMessage)
	signature := base64.StdEncoding.EncodeToString((mac.Sum(nil)))

	// берем адрес из конфига
	url := beego.AppConfig.String("apiGetParticipants")+"id="+noticeId+"&signature="+signature

	err, str, resp := curl.String(url,"timeout=10")

	if err != nil || resp.StatusCode != 200 {
		fmt.Println("User: Ошибка получения данных пользователя. Error: ", err, " HttpStatusCode: ", resp.StatusCode)
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
