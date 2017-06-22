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
func NewConn() (conn *redis.Client, err error) {

	if host == "" {
		initRedisParams()
	}

	conn = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	return
}

// init redis params
func initRedisParams() {

	host = os.Getenv("RedisHost")
	password = os.Getenv("RedisPassword")
	db, _ = strconv.Atoi(os.Getenv("RedisDb"))

	if host == "" {
		println("Redis host Is Empty")
	}

	return
}

// example
func RedisExample()  {
	conn, err := NewConn()

	println(err)

	result, err := conn.Ping().Result()

	println(result)
	println(err)
}
