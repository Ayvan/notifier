package models

/**
	Параметры каналов, включенных в сервисе
	Сюда следует добавлять новые каналы после их реализации
 */

func GetChannels() []Channel {

	//число каналов указываем тут
	channelsNum := 1

	channels := make([]Channel, channelsNum)

	//тут забиваем массив объектами каналов
	channels[0] = NewEmailChannel()
	//channels[1] = NewSmsChannel()

	return channels
}
