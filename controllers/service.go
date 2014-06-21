package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"time"
	"iforgetgo/models"
)

type ServiceController struct {
	beego.Controller
}

/**
	Читатель БД, он запрашивает в БД уведомления, которые надо отправить в ближайшее время,
	отправляет их дальше, а также решает, удаляем ли это сообщение из БД или нет, удаляемые отправляет в noticeCleanChan
 */
func (this *ServiceController) DbReader(noticeChan chan *models.Notice, noticeCleanChan chan *models.Notice) {
	ch := time.Tick(2 * time.Second)
	notice := models.Notice{1,"notice","message","20.12.2012 00:00:00", 1,1}
	for{
		select {
		case <- ch:
			fmt.Println("Read ok!!!")
			noticeChan <- &notice
			noticeCleanChan <- &notice
		}
	}
}

/**
	"Чистильщик" БД, получает из chan уведомления и удаляет их
 */
func (this *ServiceController) DbCleaner(noticeCleanChan chan *models.Notice) {
	for {
		<-noticeCleanChan
		fmt.Println("Clean ok!!!")
	}
}

/**
	Обработчик уведомлений: получает уведомление, из Group получает список пользователей и отправляет им сообщения
 */
func (this *ServiceController) NoticeWorker(noticeChan chan *models.Notice, messageChan chan *models.Message) {
	for {
		<-noticeChan
		fmt.Println("Notice worker ok!")
		message := models.Message{1,1,"message"}
		messageChan <- &message
		/**
			обрабатывает полученный notice
			запрашивает у группы список пользователей
			определяет список пользователей, кому его отправить
			отправляет в MessageWorker
	 	*/
	}
}

/**
	Обработчик сообщений: получает сообщения и ID пользователей, из User получает список каналов и их параметры
	 и отправляет сообщения	в соответствующие каналы, передавая адрес получателя (телефон, email и т.д.)
 */
func (this *ServiceController) MessageWorker(messageChan chan *models.Message, channelMessageChan chan *models.ChannelMessage) {
	for {<-messageChan
		fmt.Println("Message worker ok!")
		channelMessage := models.ChannelMessage{1,1,"message"}
		channelMessageChan <- &channelMessage
		/**
			получает пользователя (ID) кому отправить
			запрашивает у UserModel список каналов
			отправляет в ChannelDispatcher сообщение и номер канала
	 	*/
	}
}

/**
	Диспетчер каналов: он знает обо всех каналах, создает для них набор chan (по 1 на канал) и запускает воркеры,
	которые будут обрабатывать сообщения, адресованные их каналам
 */
func (this *ServiceController) ChannelDispatcher(channelMessageChan chan *models.ChannelMessage) {
	for {
		<-channelMessageChan
		fmt.Println("Channel dispatcher ok!")

		channel1 := models.NewChannelEmail()
		chan1 := make(chan *models.ChannelMessage)

		go this.ChannelMessageWorker(&channel1, chan1)

		/**
			создает chan для всех каналов (по 1 на канал)
			запускает ChannelMessageWorker'ы (от 1 до многих на каждый chan)
			роутит сообщения в каналы, прослушиваемые несколькими (от 1 на канал) ChannelMessageWorker
	 	*/
	}
}

/**
	Обработчик сообщений, отправленных в канал: получает адрес и сообщение, запускает метод Channel.Send()
	 Метод Channel.Send() должен отформатировать сообщение согласно правилам канала и вызывать соответствующий сервис-провайдер
 */
func (this *ServiceController) ChannelMessageWorker(channel models.Channel, channelMessageChan chan *models.ChannelMessage) {
	message := <- channelMessageChan
	channel.Send(message)
}
