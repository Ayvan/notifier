package models

func GetChannels() []Channel{

	//число каналов указываем тут
	channelsNum := 2

	channels := make([]Channel,channelsNum)

	//тут забиваем массив объектами каналов
	channels[0] = NewEmailChannel()
	channels[1] = NewSmsChannel()

	return channels
}
