package redis

import (
	"BlueBell/settings"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func Init(ctg *settings.RedisConfig) (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			ctg.Host,
			ctg.Port,
		),
		Password: ctg.Password,
		DB:       ctg.DB,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil

}

func Close() {
	if rdb != nil {
		rdb.Close()
	}
}
