package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"time"
	"iforgetgo/models"
	"iforgetgo/services"
)

type ServiceController struct {
	beego.Controller
}

/**
	Читатель БД, он запрашивает в БД уведомления, которые надо отправить в ближайшее время,
	отправляет их дальше, а также решает, удаляем ли это сообщение из БД или нет, удаляемые отправляет в noticeCleanChan
 */
func (this *ServiceController) DbReader(noticeChan chan *models.Notice, noticeCleanChan chan *models.Notice, redis services.Redis) {
	redis.Connect()
	ch := time.Tick(2 * time.Second)
	for {
		select {
		case <-ch:

			notices := models.NewNoticesFromRedis(redis)

			for _, notice := range notices {
				if (notice != nil) {
					noticeChan <- notice
					noticeCleanChan <- notice

					fmt.Println("Notice " + notice.Id + " pushed!")
				}
			}

			fmt.Println("DbReader finished")
		}
	}
	redis.Disconnect()
}

/**
	"Чистильщик" БД, получает из chan уведомления и удаляет их
 */
func (this *ServiceController) DbCleaner(noticeCleanChan chan *models.Notice, redis services.Redis) {
	redis.Connect()
	for {
		notice := <-noticeCleanChan
		redis.Delete(notice.Id)
		redis.DeleteFromRange("notices", notice.Id)

		fmt.Printf("Clean ok! Notice id: %s\n", notice.Id)
	}
	redis.Disconnect()
}

/**
	Обработчик уведомлений: получает уведомление, из Group получает список пользователей и отправляет им сообщения
 */
func (this *ServiceController) NoticeWorker(noticeChan chan *models.Notice, messageChan chan *models.Message, redis services.Redis) {
	redis.Connect()
	for {
		notice := <-noticeChan // читаем notice
		fmt.Println("Notice worker ok!", notice)

		// получаем группу из нотиса
		group := models.FindGroup(notice.Group, redis)
		fmt.Println("NoticeWorker group: ", group)
		//получаем список пользователей группы

		fmt.Println("NoticeWorker members: ", group.Members)
		// отправляем в MessageWorker все сообщения
		for _, member := range group.Members {
			message := models.NewMessage(notice.Id, notice.Author, member, notice.Message)
			messageChan <- message
		}
	}
	redis.Disconnect()
}

/**
	Обработчик сообщений: получает сообщение и пользователей,
	из User получает список каналов
	и отправляет сообщения	в соответствующие каналы, передавая адрес получателя (телефон, email и т.д.)
 */
func (this *ServiceController) MessageWorker(messageChan chan *models.Message, channelMessageChan chan *models.ChannelMessage, redis services.Redis) {
	redis.Connect()
	for {
		select {
		case message := <-messageChan:

			fmt.Println("MessageWorker: ", "Принял", message)

			//Каналы и адреса
			addresses := models.FindUserAddresses(message.Receiver, redis)

			//Получатель
			receiver := models.FindUser(message.Receiver)
		for _, address := range addresses {
			//Формируем сообщение для оправки в воркер каналов
			channelMessage := models.NewChannelMessage("1", address.Channel, message.Message, address.Address, receiver.Name)
			channelMessageChan <- channelMessage
			fmt.Println("MessageWorker: ", "Отправлено в очередь", channelMessage)
		}

			fmt.Println("MessageWorker: ", "Message worker ok!")
		}

		/**
			получает пользователя (ID) кому отправить
			запрашивает у UserModel список каналов
			отправляет в ChannelDispatcher сообщение и номер канала
	 	*/
	}
	redis.Disconnect()
}

/**
	Диспетчер каналов: он знает обо всех каналах, создает для них набор chan (по 1 на канал) и запускает воркеры,
	которые будут обрабатывать сообщения, адресованные их каналам
 */
func (this *ServiceController) ChannelDispatcher(channelMessageChan chan *models.ChannelMessage) {

	fmt.Println("ChannelDispatcher: Запущен")

	channels := models.GetChannels()
	chansForChannels := make([]chan *models.ChannelMessage, len(channels))

	for i, channel := range channels {
		//берем каждый канал, создаем для него chan и запускаем горутину
		chansForChannels[i] = make(chan *models.ChannelMessage)
		go this.ChannelMessageWorker(channel, chansForChannels[i])
		fmt.Println("ChannelDispatcher: ", "Создан воркер для канала ", channel.GetName())
	}

	//запускаем роутер, их может быть много
	go this.ChannelRouter(channelMessageChan, channels, chansForChannels)

	/**
		создает chan для всех каналов (по 1 на канал)
		запускает ChannelMessageWorker'ы (от 1 до многих на каждый chan)
		роутит сообщения в каналы, прослушиваемые несколькими (от 1 на канал) ChannelMessageWorker
	 */

}

func (this *ServiceController) ChannelRouter(channelMessageChan chan *models.ChannelMessage, channels []models.Channel, chansForChannels []chan *models.ChannelMessage) {
	for {
		//возьмем из очереди сообщение
		channelMessage := <-channelMessageChan
		fmt.Println("ChannelRouter: ", "Получено сообщение", channelMessage)
		//переберем все каналы
		for i, channel := range channels {
			//если канал соответствует каналу в сообщении, то отправим
			if channel.GetName() == channelMessage.Channel {
				fmt.Println("ChannelRouter: ", "Сообщение отправлено в канал", channelMessage.Channel)
				chansForChannels[i] <- channelMessage
			}
		}
	}
}

/**
	Обработчик сообщений, отправленных в канал: получает адрес и сообщение, запускает метод Channel.Send()
	 Метод Channel.Send() должен отформатировать сообщение согласно правилам канала и вызывать соответствующий сервис-провайдер
 */
func (this *ServiceController) ChannelMessageWorker(channel models.Channel, channelMessageChan chan *models.ChannelMessage) {
	for {
		channelMessage := <-channelMessageChan
		fmt.Println("ChannelMessageWorker: ", "Сообщение отправлено в канал", channelMessage)
		channel.Send(channelMessage)
	}
}
