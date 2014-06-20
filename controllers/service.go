package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"time"
	"reminder/models"
)

type ServiceController struct {
	beego.Controller
}

func (this *ServiceController) DbReader(noticeChan chan models.Notice, noticeCleanChan chan models.Notice) {
	ch := time.Tick(2 * time.Second)
	notice := models.Notice{1}
	for{
		select {
		case <- ch:
			fmt.Println("Read ok!!!")
			noticeChan <- notice
			//noticeCleanChan <- &notice
		}
	}
}

func (this *ServiceController) DbCleaner(noticeCleanChan chan models.Notice) {
	//<- noticeCleanChan
}

func (this *ServiceController) NoticeWorker(noticeChan chan models.Notice, messageChan chan models.Message) {
	<- noticeChan
	fmt.Println("Notice worker ok!")
	message := models.Message{1}
	messageChan <- message
	/**
	обрабатывает полученный notice
	запрашивает у канала список пользователей
	определяет список пользователей, кому его отправить
	отправляет в MessageWorker
	 */

}

func (this *ServiceController) MessageWorker(messageChan chan models.Message, channelMessageChan chan models.ChannelMessage) {
	<- messageChan
	fmt.Println("Message worker ok!")
	channelMessage := models.ChannelMessage{1}
	channelMessageChan <- channelMessage
	/**
	получает пользователя (ID) кому отправить
	запрашивает у UserModel список каналов
	отправляет в ChannelDispatcher сообщение и номер канала
	 */
}

func (this *ServiceController) ChannelDispatcher(channelMessageChan chan models.ChannelMessage) {
	<- channelMessageChan
	fmt.Println("Channel dispatcher ok!")
	/**
	создает chan для всех каналов (по 1 на канал)
	запускает ChannelMessageWorker'ы (от 1 до многих на каждый chan)
	роутит сообщения в каналы, прослушиваемые несколькими (от 1 на канал) ChannelMessageWorker
	 */
}

func (this *ServiceController) ChannelMessageWorker() {

}
