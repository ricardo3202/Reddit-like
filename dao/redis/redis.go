package redis

import (
	"bluebell/settings"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

//Addr: fmt.Sprintf("%s:%d",
//viper.GetString("redis.host"),
//viper.GetInt("redis.port"),
//),

func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.Db,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// Background返回一个非空的Context。它永远不会被取消，没有值，也没有截止日期。
	// 它通常由main函数、初始化和测试使用，并作为传入请求的顶级上下文

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Error("connect redis failed", zap.Error(err))
		return
	}
	return
}

func Close() {
	_ = client.Close()
}
