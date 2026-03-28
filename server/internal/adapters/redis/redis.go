package redis

import (
	"context"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addr, username, password string) *goredis.Client {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Printf("connected to redis at %s", addr)

	return rdb
}
