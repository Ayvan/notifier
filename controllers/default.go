package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"notifier/models"
	"notifier/services"
)

type MainController struct {
	beego.Controller
}

type TestController struct {
	beego.Controller
}

func (this *TestController) Get() {

	users := models.GetParticipants("1")

	fmt.Printf("%v", users)

	this.Data["Website"] = "";
	this.Data["Email"] = "";
	this.TplNames = "index.tpl"
}


func (this *MainController) Get() {

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
	// объявим структуру, отвечающую кодом 0 и сообщением "успешно"
	noticeId := this.GetString("id");
	noticeMessage := this.GetString("message")
	noticeTime := this.GetString("noticeTime")

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

func (this *MainController) deleteNotice(redis services.Redis) {
	noticeId := this.GetString("id");

	redis.Delete("notice:"+noticeId)
	redis.DeleteFromRange("notices","notice:"+noticeId)

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
	noticeTime := this.GetString("noticeTime")

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

