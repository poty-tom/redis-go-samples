package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	conn *redis.Client
)

func init() {
	// redisクライアントを初期化
	conn = redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})
}

func main() {
	ctx := context.Background()

	// セッション情報のセット
	err := conn.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		fmt.Println("error occurred when set key/value")
	} else {
		fmt.Println("set key value")
	}

	// セッション情報の取得
	v, err := conn.Get(ctx, "key").Result()
	if err == redis.Nil {
		fmt.Println("redis key not exists")
	} else if err != nil {
		fmt.Println(err)
	} else if v == "" {
		fmt.Println("redis value is empty")
	} else {
		fmt.Printf("redis value is:%v\n", v)
	}
}
