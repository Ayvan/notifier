package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"notifier/services"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

type MainController struct {
	beego.Controller
}

type TestController struct {
	beego.Controller
}


func (this *MainController) Get() {
	response := struct {
			Code    int
			Message string
		} {
		-3, "only POST requests are allowed",
	}

	this.Data["json"] = &response;
	this.ServeJson()
}

func (this *MainController) Post() {

	// запускаем проверку подписи (ЭЦП)
	if this.checkSignature() != true {

		response := struct {
				Code    int
				Message string
			} {
			-2, "signature error",
		}

		this.PrintDevLn("MainController: Signature error!")

		this.Data["json"] = &response;
		this.ServeJson()
		return
	}

	this.PrintDevLn("MainController: Signature check success!")

	// получим параметр "действией", он указывает, что требуется сделать
	action := this.GetString("action");

	redis := services.NewRedis(beego.AppConfig.String("redisHost"), beego.AppConfig.String("redisPort"))

	redis.Connect()

	switch action {
	case "add":
		this.addNotice(redis)
	case "delete":
		this.deleteNotice(redis)
	case "edit":
		this.editNotice(redis)
	}

	// если ни одно из действий не подошло, ответим что действие не найдено
	response := struct {
			Code    int
			Message string
		} {
		-1, "bad request: action not found",
	}

	this.Data["json"] = &response;
	this.ServeJson()
}

func (this *MainController) addNotice(redis services.Redis) {
	this.editNotice(redis)
}

func (this *MainController) deleteNotice(redis services.Redis) {
	noticeId := this.GetString("id");

	redis.Delete("notice:"+noticeId)
	redis.DeleteFromRange("notices","notice:"+noticeId)

	this.PrintDevLn("MainController: удалено событие "+noticeId)

	response := struct {
			Code    int
			Message string
		} {
		0, "success",
	}

	this.Data["json"] = &response;
	this.ServeJson()
}

func (this *MainController) editNotice(redis services.Redis) {
	noticeId := this.GetString("id");
	noticeMessage := this.GetString("message")
	noticeTime := this.GetString("time")

	result := redis.Get("notice:"+noticeId)

	if len(result) >=1 {
		this.PrintDevLn("MainController: отредактировано событие "+noticeId)
	} else {
		this.PrintDevLn("MainController: новое событие "+noticeId)
	}

	// удаляем старое уведомление из очереди и из списка уведомлений
	redis.Delete("notice:"+noticeId)
	redis.DeleteFromRange("notices","notice:"+noticeId)

	// добавляем новое уведомление
	redis.AddNotice(noticeId, noticeMessage, noticeTime)

	response := struct {
			Code    int
			Message string
		} {
		0, "success",
	}

	this.Data["json"] = &response;
	this.ServeJson()
}

func (this *MainController) checkSignature() bool{
	action := this.GetString("action");
	noticeId := this.GetString("id");
	noticeMessage := this.GetString("message")
	noticeTime := this.GetString("time")
	signature := this.GetString("signature")

	key := []byte(beego.AppConfig.String("apiSecretKey"))
	mac := hmac.New(sha1.New,key)
	macMessage := []byte(action+noticeId+noticeMessage+noticeTime)
	mac.Write(macMessage)
	expectedMAC := mac.Sum(nil)
	signatureMAC, _ := base64.StdEncoding.DecodeString(signature)
	isMacEquals := hmac.Equal(signatureMAC,expectedMAC)

	return isMacEquals
}

/**
	Вывод сообщений для разработки, выводит в консоль сообщения только если включен параметр devtrace
*/
func (this *MainController) PrintDevLn(a ...interface{}) {
	_ = a
	devtrace, _ := beego.AppConfig.Bool("devtrace")
	if devtrace {
		fmt.Println(a...)
	}
}
