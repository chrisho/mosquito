package redis

import (
	"os"
	"strconv"
	"github.com/go-redis/redis"
)

var (
	host     string
	password string
	db       int
)

// create redis connection
func NewConn() (client *redis.Client, err error) {

	if host == "" {
		initRedisParams()
	}

	client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	return
}

// init redis parameter
func initRedisParams() {

	host = os.Getenv("RedisHost")
	password = os.Getenv("RedisPassword")
	db, _ = strconv.Atoi(os.Getenv("RedisDb"))

	if host == "" {
		println("Redis host Is Empty")
	}

	return
}
