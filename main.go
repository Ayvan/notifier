package main

import (
	_ "iforgetgo/routers"
	"github.com/astaxie/beego"
	"iforgetgo/controllers"
	"iforgetgo/models"
)

func startService() {
	c := controllers.ServiceController{}
	// for i:=0;i<N;i++ { запуск нескольких горутин воркеров
	noticeChan := make(chan models.Notice, 100)
	noticeCleanChan := make(chan models.Notice,	100)
	messageChan := make(chan models.Message, 100)
	channelMessageChan := make(chan models.ChannelMessage, 100)

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

