package main

import (
	"fmt"
	"os"
	"os/signal"
	_ "iforgetgo/routers"
	"github.com/astaxie/beego"
	"iforgetgo/controllers"
	"iforgetgo/models"
	"iforgetgo/services"
)

func startService() {

	c := controllers.ServiceController{}
	c.InitService()
	// for i:=0;i<N;i++ { запуск нескольких горутин воркеров

	/******************************************Создание каналов*******************************************************/
	/**
		Поступает информация о текущей нотификации
		Поля - группа, автор, текст сообщения
	 */
	noticeChan := make(chan *models.Notice, 100)

	/**
		Поступает инфомация о нотификации для ее удаления
	 */
	noticeCleanChan := make(chan *models.Notice, 100)

	/**
		Поступает информация о сообщении для конкретного пользователя
		Поля - получатель, отправитель, сообщение
	 */
	messageChan := make(chan *models.Message, 100)

	/**
		Поступает информация для отправки сообщения в конкретный канал
		Поля - получатель, сообщение, название канала, имя получателя
	 */
	channelMessageChan := make(chan *models.ChannelMessage, 100)

	// Подключаемся к redis
	redis := services.NewRedis(beego.AppConfig.String("redisHost"), beego.AppConfig.String("redisPort"))

	/******************************************Создание процессов******************************************************/

	//запускаем процесс, читающий БД
	go c.DbReader(noticeChan, noticeCleanChan, redis)

	//запускаем процесс, удаляющий из БД обработанные записи
	go c.DbCleaner(noticeCleanChan, redis)

	//запускаем воркер нотификаций - выбирает получателей из группы для отправки им сообщений
	go c.NoticeWorker(noticeChan, messageChan, redis)

	//запусукаем воркер сообщений - выбирает каналы пользователя, в которые отправлять сообщение
	go c.MessageWorker(messageChan, channelMessageChan, redis)

	//запусаем диспетчер каналов - создает chan для каждого канала и воркеры для обработки этих chan
	go c.ChannelDispatcher(channelMessageChan)

	// останов сервера по Ctrl-C
	<-sigChan
	c.Stop()
	fmt.Println("Сервер успешно остановлен.")
	os.Exit(0)
}

func main() {
	startService()
	beego.Run()
}

