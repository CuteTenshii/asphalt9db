package utils

import "github.com/redis/go-redis/v9"

var rdb *redis.Client

func GetRedis() *redis.Client {
	if rdb != nil {
		return rdb
	}

	rdb = redis.NewClient(&redis.Options{})
	return rdb
}
