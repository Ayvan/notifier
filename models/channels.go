package models

func GetChannels() []Channel{

	channels := make([]Channel,1)

	channels[0] = NewChannelEmail()

	return channels
}
