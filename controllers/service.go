package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"notifier/models"
	"notifier/services"
	"sync"
	"time"
)

type ServiceController struct {
	wg               *sync.WaitGroup
	quitChanDbReader chan bool
	finishedDbReader chan bool

	quitChanDbCleaner chan bool
	finishedDbCleaner chan bool

	quitChanNoticeWorker chan bool
	finishedNoticeWorker chan bool

	quitChanMessageWorker chan bool
	finishedMessageWorker chan bool

	quitChanChannelDispatcher chan bool
	finishedChannelDispatcher chan bool

	quitChanChannelRouter chan bool
	finishedChannelRouter chan bool

	beego.Controller
}

/**
	Инициализатор сервиса, создает каналы
 */
func (this *ServiceController) InitService() {
	this.wg = &sync.WaitGroup{}

	//каналы для команд завершения
	this.quitChanDbReader = make(chan bool, 1)
	this.quitChanDbCleaner = make(chan bool, 1)
	this.quitChanNoticeWorker = make(chan bool, 1)
	this.quitChanMessageWorker = make(chan bool, 1)
	this.quitChanChannelDispatcher = make(chan bool, 1)
	this.quitChanChannelRouter = make(chan bool, 1)

	//каналы для отчетов об успешном завершении
	this.finishedDbReader = make(chan bool, 1)
	this.finishedDbCleaner = make(chan bool, 1)
	this.finishedNoticeWorker = make(chan bool, 1)
	this.finishedMessageWorker = make(chan bool, 1)
	this.finishedChannelDispatcher = make(chan bool, 1)
	this.finishedChannelRouter = make(chan bool, 1)
}
/**
	Запуск сервиса
 */
func (this *ServiceController) Run() {
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

	// подключаемся к redis
	redis := services.NewRedis(beego.AppConfig.String("redisHost"), beego.AppConfig.String("redisPort"))

	/******************************************Создание процессов******************************************************/

	//запускаем процесс, читающий БД
	go this.DbReader(noticeChan, noticeCleanChan, redis)

	//запускаем процесс, удаляющий из БД обработанные записи
	go this.DbCleaner(noticeCleanChan, redis)

	//запускаем воркер нотификаций - запрашивает у системы-пользователя получателей для отправки им сообщений
	go this.NoticeWorker(noticeChan, messageChan, redis)

	//запусукаем воркер сообщений - выбирает каналы пользователя, в которые отправлять сообщение
	go this.MessageWorker(messageChan, channelMessageChan, redis)

	//запусаем диспетчер каналов - создает chan для каждого канала и воркеры для обработки этих chan
	go this.ChannelDispatcher(channelMessageChan)

}

/**
Читатель БД, он запрашивает в БД уведомления, которые надо отправить в ближайшее время,
отправляет их в канал обработки уведомлений и канал удаления
*/
func (this *ServiceController) DbReader(noticeChan chan *models.Notice, noticeCleanChan chan *models.Notice, redis services.Redis) {
	this.wg.Add(1)
	fmt.Println("DbReader: STARTED")
	redis.Connect()
	defer func() {
		redis.Disconnect()
		this.wg.Done()
		fmt.Println("DbReader: STOPPED")
		this.finishedDbReader <- true
	}()
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

				this.PrintDevLn("DbReader: ", "Уведомление ", notice.Id, " отправлено в обработку и удаление")
			}
		}

			fmt.Println("DbReader: ", "Обработка закончена")

		case <-this.quitChanDbReader:
			return

		}
	}
}

/**
"Чистильщик" БД, получает из chan уведомления и удаляет их
*/
func (this *ServiceController) DbCleaner(noticeCleanChan chan *models.Notice, redis services.Redis) {
	this.wg.Add(1)
	fmt.Println("DbCleaner: STARTED")
	redis.Connect()
	defer func() {
		redis.Disconnect()
		this.wg.Done()
		fmt.Println("DbCleaner: STOPPED")
		this.finishedDbCleaner <- true
	}()

	for {
		select {
		case notice := <-noticeCleanChan:
			redis.DeleteFromRange("notices", "notice:"+notice.Id)
			this.PrintDevLn("DbCleaner: ", "Уведомление ", notice.Id, " удалено")
		case <-this.quitChanDbCleaner:
			return
		}
	}
}

/**
Обработчик уведомлений: получает уведомление, запрашивает список пользователей и отправляет в обработчик сообщений
*/
func (this *ServiceController) NoticeWorker(noticeChan chan *models.Notice, messageChan chan *models.Message, redis services.Redis) {
	this.wg.Add(1)
	fmt.Println("NoticeWorker: STARTED")
	redis.Connect()
	defer func() {
		redis.Disconnect()
		this.wg.Done()
		fmt.Println("NoticeWorker: STOPPED")
		this.finishedNoticeWorker <- true
	}()
	for {
		select {

		case notice := <-noticeChan: // читаем notice
			this.PrintDevLn("NoticeWorker: ", "Обработка уведомления ", notice.Id)
			continue
			// получаем список пользователей, подписанных на событие
			users := models.GetParticipants(notice.Id)

			if users == nil {
				continue
			}

			this.PrintDevLn("NoticeWorker: ", "Найдены группа пользователей для оповещения ", users.NoticeId)

			//получаем список пользователей группы
			this.PrintDevLn("NoticeWorker: ", "В группе найдено ", len(users.Users), " получателей")

			// отправляем в MessageWorker все сообщения
		for _, member := range users.Users {
			message := models.NewMessage(notice.Id, notice.Message, member)
			messageChan <- message

			this.PrintDevLn("NoticeWorker: ", "Сообщение для пользователя ", member, " отправлено")
		}

			this.PrintDevLn("NoticeWorker: ", "Закончил обработку уведоления")

		case <-this.quitChanNoticeWorker:
			return

		}
	}
}

/**
Обработчик сообщений: получает сообщение и получателя,
Отправляет сообщения в обработчик каналов, передавая адрес получателя, канал, имя получателя, текст сообщения
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

			this.PrintDevLn("MessageWorker: ", "Принял сообщение")
			//Каналы и адреса
			addresses := message.Receiver.Addresses;
			this.PrintDevLn("MessageWorker: ", "Найдено ", len(addresses), " каналов")

			if addresses == nil {
				this.PrintDevLn("MessageWorker: ", "Не найдены адреса! "+message.Receiver.Name)
				continue
			}
			this.PrintDevLn("MessageWorker: ", "Найден получатель "+message.Receiver.Name)

		for _, address := range addresses {
			//Формируем сообщение для оправки в воркер каналов
			channelMessage := models.NewChannelMessage(message.NoticeId, address.Channel, message.Message, address.Address, message.Receiver.Name)
			channelMessageChan <- channelMessage
			this.PrintDevLn("MessageWorker: ", "Отправлено в очередь", channelMessage)
		}

			this.PrintDevLn("MessageWorker: ", "Message worker ok!", message.Receiver.Name)

		case <-this.quitChanMessageWorker:
		this.finishedMessageWorker <- true
			return

		}
	}

}

/**
Диспетчер каналов: он знает обо всех каналах, создает для них набор chan (по 1 на канал) и запускает воркеры,
которые будут обрабатывать сообщения, адресованные их каналам
*/
func (this *ServiceController) ChannelDispatcher(channelMessageChan chan *models.ChannelMessage) {

	fmt.Println("ChannelDispatcher: STARTED")
	this.wg.Add(1)
	defer func() {
		this.wg.Done()
		fmt.Println("ChannelDispatcher: STOPPED")
		this.finishedChannelDispatcher <- true
	}()

	channels := models.GetChannels()
	chansForChannels := make([]chan *models.ChannelMessage, len(channels))

	quitChansForChannels := make([]chan bool, len(channels))
	finishChansForChannels := make([]chan bool, len(channels))

	for i, channel := range channels {
		//берем каждый канал, создаем для него chan и запускаем горутину
		chansForChannels[i] = make(chan *models.ChannelMessage)
		quitChansForChannels[i] = make(chan bool)
		finishChansForChannels[i] = make(chan bool)
		go this.ChannelMessageWorker(channel, chansForChannels[i], quitChansForChannels[i], finishChansForChannels[i])
		this.PrintDevLn("ChannelDispatcher: ", "Создан воркер для канала ", channel.GetName())
	}

	//запускаем роутеры
	go this.ChannelRouter(channelMessageChan, channels, chansForChannels)

	//ожидаем команды отключения диспетчера
	<-this.quitChanChannelDispatcher

	//останавливаем роутер и ждем его завершения
	this.quitChanChannelRouter <- true
	<-this.finishedChannelRouter

	//останавливаем каналы и ждем завершение каждого
	for i, _ := range channels {
		quitChansForChannels[i] <- true
		<-finishChansForChannels[i]
	}

}

func (this *ServiceController) ChannelRouter(channelMessageChan chan *models.ChannelMessage, channels []models.Channel, chansForChannels []chan *models.ChannelMessage) {
	fmt.Println("ChannelRouter: STARTED")
	this.wg.Add(1)
	defer func() {
		this.wg.Done()
		fmt.Println("ChannelRouter: STOPPED")
		this.finishedChannelRouter <- true
	}()

	for {
		select {
		case channelMessage := <-channelMessageChan:
			//возьмем из очереди сообщение
			this.PrintDevLn("ChannelRouter: ", "Получено сообщение", channelMessage)
			//переберем все каналы
		for i, channel := range channels {
			//если канал соответствует каналу в сообщении, то отправим
			if channel.GetName() == channelMessage.Channel {
				chansForChannels[i] <- channelMessage
				this.PrintDevLn("ChannelRouter: ", "Сообщение отправлено в канал", channelMessage.Channel)
			}
		}

		case <-this.quitChanChannelRouter:
			return
		}
	}
}

/**
Обработчик сообщений, отправленных в канал: получает адрес и сообщение, запускает метод Channel.Send()
 Метод Channel.Send() должен отформатировать сообщение согласно правилам канала и вызывать соответствующий сервис-провайдер
*/
func (this *ServiceController) ChannelMessageWorker(channel models.Channel, channelMessageChan chan *models.ChannelMessage, quitChan, finishChan chan bool) {
	fmt.Println("ChannelMessageWorker: STARTED ", channel.GetName())
	this.wg.Add(1)
	defer func() {
		this.wg.Done()
		fmt.Println("ChannelMessageWorker: STOPPED", channel.GetName())
		finishChan <- true
	}()

	for {
		select {

		case channelMessage := <-channelMessageChan:
			this.PrintDevLn("ChannelMessageWorker: ", "Сообщение отправлено в канал", channelMessage)
			channel.Send(channelMessage)

		case <-quitChan:
			return
		}
	}
}

/**
	Последовательная остановка параллельно работающих служб сервиса.
	Каждой службе по-очереди отправляется сообщение, что требуется завершить работу,
	служба заканчивает текущую работу (без потери данных) и отвечает, что работа завершена,
	затем команда завершения передается следующей службе.
	Последовательность остановки служб сервиса сделана таким образом, что остановка происходит в порядке работы
	конвейера данных.
 */
func (this *ServiceController) Stop() {

	this.quitChanDbReader <- true
	<-this.finishedDbReader

	this.quitChanDbCleaner <- true
	<-this.finishedDbCleaner

	this.quitChanNoticeWorker <- true
	<-this.finishedNoticeWorker

	this.quitChanMessageWorker <- true
	<-this.finishedMessageWorker

	this.quitChanChannelDispatcher <- true
	<-this.finishedChannelDispatcher

	this.wg.Wait()
}

/**
	Вывод сообщений для разработки, выводит в консоль сообщения только если включен параметр devtrace
*/
func (this *ServiceController) PrintDevLn(a ...interface{}) {
	_ = a
	devtrace, _ := beego.AppConfig.Bool("devtrace")
	if devtrace {
		fmt.Println(a...)
	}
}
