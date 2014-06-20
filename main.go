package main

import (
	_ "reminder/routers"
	"github.com/astaxie/beego"
	"reminder/controllers"
	"reminder/models"
)

func startService() {
	c := controllers.ServiceController{}
	// for i:=0;i<N;i++ { запуск нескольких горутин воркеров
	noticeChan := make(chan models.Notice)
	noticeCleanChan := make(chan models.Notice)
	messageChan := make(chan models.Message)
	channelMessageChan := make(chan models.ChannelMessage)

	go c.DbReader(noticeChan,noticeCleanChan)
	go c.DbCleaner(noticeCleanChan)
	go c.NoticeWorker(noticeChan,messageChan)
	go c.MessageWorker(messageChan, channelMessageChan)
	go c.ChannelDispatcher(channelMessageChan)
}

func main() {
	startService() //
	beego.Run()
}

