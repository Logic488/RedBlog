package redis

import (
	"RedBlog/settings"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 声明一个全局的rdb变量
var client *redis.Client

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		//Addr:     "localhost:6379",
		//Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.addr"), viper.GetInt("redis.port")),
		//Password: viper.GetString("redis.password"), // no password set
		//DB:       viper.GetInt("redis.db"),          // use default DB
		//PoolSize: viper.GetInt("redis.pool_size"),
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping(context.Background()).Result()

	return
}

func Close() {
	_ = client.Close()
}
