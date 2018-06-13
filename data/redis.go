package data

import (
	"strconv"
	"easy-go/conf"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

// this style has a weak point,
// when app begin running this init will be running too whatever you indeed need
func init() {
	db, err := strconv.Atoi(conf.AppConfig["redisdb"])
	if err != nil {
		db = 0 // db 默认设置0
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.AppConfig["redisaddr"],
		Password: conf.AppConfig["redispassword"],
		DB:       db,
	})
}
