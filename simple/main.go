package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	set, err := rdb.SetNX(ctx, "key", "value", 0).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(set)

	val, err := rdb.Get(ctx, "key-ring3").Result()
	switch {
	case err == redis.Nil:
		fmt.Println("key does not exist")
	case err != nil:
		fmt.Println("Get failed", err)
	case val == "":
		fmt.Println("value is empty")
	}
	fmt.Println(val)
}
