package redis

import (
	"github.com/go-redis/redis"
	"os"
	"strconv"
)

var (
	host     string
	password string
	db       int
)

func GetRedisClient() (client *redis.Client) {
	return
}

// create redis connection
func NewConn() (conn *redis.Client, err error) {

	conn, err = NewConnDB(0)

	return
}

func NewConnDB(select_db int) (conn *redis.Client, err error) {

	if host == "" {
		initRedisParams()
	}

	conn = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       select_db,
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
		host = "127.0.0.1:6379"
	}

	return
}
