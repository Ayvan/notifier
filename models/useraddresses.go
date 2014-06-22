package models

import (
	"iforgetgo/services"
)

type UserAddresses struct {
	Channel string
	Address string
}

func NewUserAddresses(channel string, address string) *UserAddresses {
	return &UserAddresses{channel, address}
}

func FindUserAddresses(uuid string, redis services.Redis) []*UserAddresses {

	result := redis.Get(uuid + ":addresses")

	addresses := make([]*UserAddresses, len(result)/2, len(result)/2)
	j := 0
	for i := 0; i < len(result); i++ {
		key := i
		i = i + 1
		value := i

		channel := result[key]
		address := result[value]

		addresses[j] = NewUserAddresses(channel, address)
		j++
	}

	return addresses
}
