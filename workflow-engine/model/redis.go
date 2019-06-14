package model

import (
	"fmt"
	"time"

	redis "github.com/go-redis/redis"
)

// RedisClient redis客户端
var RedisClient *redis.Client

// SetRedis 设置redis
func SetRedis() {
	var err error
	fmt.Println("-------启动redis--------")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost + ":" + conf.RedisPort,
		Password: conf.RedisPassword,
	})
	pong, err := RedisClient.Ping().Result()
	fmt.Println(pong, err)
}

// RedisSetVal 将值保存到redis
func RedisSetVal(key, value string, expiration time.Duration) error {
	return RedisClient.Set(key, value, expiration).Err()
}

// RedisGetVal 从redis获取值
func RedisGetVal(key string) (string, error) {
	return RedisClient.Get(key).Result()
}
