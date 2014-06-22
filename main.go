package main

import (
	_ "iforgetgo/routers"
	"github.com/astaxie/beego"
	"iforgetgo/controllers"
	"iforgetgo/models"
	"iforgetgo/services"
)

func startService() {
	c := controllers.ServiceController{}
	// for i:=0;i<N;i++ { запуск нескольких горутин воркеров
	noticeChan := make(chan *models.Notice, 100)
	noticeCleanChan := make(chan *models.Notice, 100)
	messageChan := make(chan *models.Message, 100)
	channelMessageChan := make(chan *models.ChannelMessage, 100)

	// подключаемся к redis
	redis := services.NewRedis(beego.AppConfig.String("redisHost"), beego.AppConfig.String("redisPort"))
	redis.Connect()

	//запускаем процесс, читающий БД
	go c.DbReader(noticeChan, noticeCleanChan, redis)
	//запускаем процесс, удаляющий из БД обработанные записи
	go c.DbCleaner(noticeCleanChan, redis)
	//запускаем воркер уведомлений: он обрабатывает уведомление и решает кому его отправить
	go c.NoticeWorker(noticeChan, messageChan)
	//запусукаем воркер сообщений: он получает сообщение и ID получателя (юзера)
	//запрашивает у User список каналов и отправляет диспетчеру каналов сообщение и идентификатор канала
	go c.MessageWorker(messageChan, channelMessageChan, redis)
	//запусаем диспетчер каналов
	//создает chan для каждого канала и воркеры для обработки этих chan
	go c.ChannelDispatcher(channelMessageChan)
}

func main() {
	startService() //
	beego.Run()
}

