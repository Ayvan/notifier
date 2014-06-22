package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"iforgetgo/models"
	"iforgetgo/services"
	"sync"
	"time"
)

type ServiceController struct {
	wg       *sync.WaitGroup
	quitChan chan bool
	beego.Controller
}

func (this *ServiceController) InitService() {
	this.wg = &sync.WaitGroup{}
	this.quitChan = make(chan bool, 1)
}

/**
Читатель БД, он запрашивает в БД уведомления, которые надо отправить в ближайшее время,
отправляет их в канал обработки уведомлений и канал удаления
*/
func (this *ServiceController) DbReader(noticeChan chan *models.Notice, noticeCleanChan chan *models.Notice, redis services.Redis) {
	fmt.Println("DbReader: ", "Запущен")
	redis.Connect()
	//Чтение из БД каждые 2 секунды
	ch := time.Tick(2 * time.Second)
	for {
		select {

		case <-ch:

			//Получаем из базы все уведомления, которые нужно сейчас отправить
			notices := models.NewNoticesFromRedis(redis)

			fmt.Println("DbReader: ", "Найдено", len(notices), "уведомлений для обработки")

			//Обходим каждое уведомление и отправляем его в канал обработки и в канал удаления
			for _, notice := range notices {
				if notice != nil {
					//Отправка в канал обработки уведомлений
					noticeChan <- notice

					//Отправка в канал удаления уведомлений
					noticeCleanChan <- notice

					fmt.Println("DbReader: ", "Уведомление ", notice.Id, " отправлено в обработку и удаление")
				}
			}

			fmt.Println("DbReader: ", "Обработка закончена")

		case <-this.quitChan:
			this.quitChan <- true
			return

		}
	}
}

/**
"Чистильщик" БД, получает из chan уведомления и удаляет их
*/
func (this *ServiceController) DbCleaner(noticeCleanChan chan *models.Notice, redis services.Redis) {
	fmt.Println("DbCleaner: ", "Запущен")
	redis.Connect()
	defer func() {
		redis.Disconnect()
		this.wg.Done()
		fmt.Println("DbCleaner: STOPPED")
	}()

	for {
		select {

		case notice := <-noticeCleanChan:
			redis.Delete(notice.Id)
			redis.DeleteFromRange("notices", notice.Id)
			fmt.Println("DbCleaner: ", "Уведомление ", notice.Id, " удалено")

		case <-this.quitChan:
			this.quitChan <- true
			return
		}
	}
}

/**
Обработчик уведомлений: получает уведомление, из Group получает список пользователей и отправляет в обработчик сообщений
*/
func (this *ServiceController) NoticeWorker(noticeChan chan *models.Notice, messageChan chan *models.Message, redis services.Redis) {
	fmt.Println("NoticeWorker: ", "Запущен")
	redis.Connect()
	defer func() {
		redis.Disconnect()
		this.wg.Done()
		fmt.Println("NoticeWorker: STOPPED")
	}()
	for {
		select {

		case notice := <-noticeChan: // читаем notice
			fmt.Println("NoticeWorker: ", "Обработка уведомления ", notice.Id)

			// получаем группу из нотиса
			group := models.FindGroup(notice.Group, redis)
			fmt.Println("NoticeWorker: ", "Найдена группа ", group.Id)

			if group == nil {
				continue
			}

			//получаем список пользователей группы
			fmt.Println("NoticeWorker: ", "В группе найдено ", len(group.Members), " получателей")

			// отправляем в MessageWorker все сообщения
			for _, member := range group.Members {
				message := models.NewMessage(notice.Id, notice.Author, member, notice.Message)
				messageChan <- message

				fmt.Println("NoticeWorker: ", "Сообщение для пользователя ", member, " отправлено")
			}

			fmt.Println("NoticeWorker: ", "Закончил обработку уведоления")

		case <-this.quitChan:
			this.quitChan <- true
			return

		}
	}
}

/**
Обработчик сообщений: получает сообщение и получателя,
Получает список каналов и адресов каналов для получателя
и отправляет сообщения в обработчик каналов, передавая адрес получателя, канал, имя получателя, текст сообщения
*/
func (this *ServiceController) MessageWorker(messageChan chan *models.Message, channelMessageChan chan *models.ChannelMessage, redis services.Redis) {
	this.wg.Add(1)
	fmt.Println("MessageWorker: STARTED")
	redis.Connect()
	defer func() {
		redis.Disconnect()
		this.wg.Done()
		fmt.Println("MessageWorker: STOPPED")
	}()
	for {
		select {

		case message := <-messageChan:

			fmt.Println("MessageWorker: ", "Принял сообщение")
			//Каналы и адреса
			addresses := models.FindUserAddresses(message.Receiver, redis)
			fmt.Println("MessageWorker: ", "Найдено ", len(addresses), " каналов")
			//Получатель
			receiver := models.FindUser(message.Receiver, redis)
			fmt.Println("MessageWorker: ", "Найден получатель "+message.Receiver)

			for _, address := range addresses {
				//Формируем сообщение для оправки в воркер каналов
				channelMessage := models.NewChannelMessage("1", address.Channel, message.Message, address.Address, receiver.Name)
				channelMessageChan <- channelMessage
				fmt.Println("MessageWorker: ", "Отправлено в очередь", channelMessage)
			}

			fmt.Println("MessageWorker: ", "Message worker ok!", receiver)

		case <-this.quitChan:
			this.quitChan <- true
			return

		}

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

	fmt.Println("ChannelDispatcher: STARTED")

	channels := models.GetChannels()
	chansForChannels := make([]chan *models.ChannelMessage, len(channels))

	for i, channel := range channels {
		//берем каждый канал, создаем для него chan и запускаем горутину
		chansForChannels[i] = make(chan *models.ChannelMessage)
		go this.ChannelMessageWorker(channel, chansForChannels[i])
		fmt.Println("ChannelDispatcher: ", "Создан воркер для канала ", channel.GetName())
	}

	//запускаем роутеры
	go this.ChannelRouter(channelMessageChan, channels, chansForChannels)

	<-this.quitChan
	fmt.Println("ChannelDispatcher: STOPPED")
	this.quitChan <- true

}

func (this *ServiceController) ChannelRouter(channelMessageChan chan *models.ChannelMessage, channels []models.Channel, chansForChannels []chan *models.ChannelMessage) {
	fmt.Println("ChannelRouter: STARTED")
	this.wg.Add(1)
	defer func() {
		this.wg.Done()
		fmt.Println("ChannelRouter: STOPPED")
	}()

	for {
		select {
		case channelMessage := <-channelMessageChan:
			//возьмем из очереди сообщение
			fmt.Println("ChannelRouter: ", "Получено сообщение", channelMessage)
			//переберем все каналы
			for i, channel := range channels {
				//если канал соответствует каналу в сообщении, то отправим
				if channel.GetName() == channelMessage.Channel {
					chansForChannels[i] <- channelMessage
					fmt.Println("ChannelRouter: ", "Сообщение отправлено в канал", channelMessage.Channel)
				}
			}

		case <-this.quitChan:
			this.quitChan <- true
			return
		}
	}
}

/**
Обработчик сообщений, отправленных в канал: получает адрес и сообщение, запускает метод Channel.Send()
 Метод Channel.Send() должен отформатировать сообщение согласно правилам канала и вызывать соответствующий сервис-провайдер
*/
func (this *ServiceController) ChannelMessageWorker(channel models.Channel, channelMessageChan chan *models.ChannelMessage) {
	fmt.Println("ChannelMessageWorker: STARTED ", channel.GetName())
	this.wg.Add(1)
	defer func() {
		this.wg.Done()
		fmt.Println("ChannelMessageWorker: STOPPED", channel.GetName())
	}()

	for {
		select {

		case channelMessage := <-channelMessageChan:
			fmt.Println("ChannelMessageWorker: ", "Сообщение отправлено в канал", channelMessage)
			channel.Send(channelMessage)

		case <-this.quitChan:
			this.quitChan <- true
			return
		}
	}
}

func (this *ServiceController) Stop() {
	this.quitChan <- true
	this.wg.Wait()
}
