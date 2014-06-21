package services

import "github.com/garyburd/redigo/redis"

type RedisConnector struct {
	Network  string;
	Host     string;
	Port     int;
}

func connect(this *RedisConnector) {
	redis.Dial(this.Network, this.Host+":"+string(this.Port));
}
