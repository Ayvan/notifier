package models

func GetChannels() []Channel{

	//число каналов указываем тут
	channelsNum := 1

	channels := make([]Channel,channelsNum)

	//тут забиваем массив объектами каналов
	channels[0] = NewChannelEmail()

	return channels
}
