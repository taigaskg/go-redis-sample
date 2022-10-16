package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"shard1": "localhost:6379",
			"shard2": "localhost:6380",
		},
	})

	ctx := context.Background()
	set, err := rdb.SetNX(ctx, "key-ring", "value-ring", 0).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(set)

	val, err := rdb.Get(ctx, "key-ring").Result()
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
