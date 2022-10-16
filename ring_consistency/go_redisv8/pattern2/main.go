package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

func main() {
	const count = 10

	fmt.Println("--------- deleting -----------")
	rdb := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"shard1": "localhost:6379",
			"shard2": "localhost:6380",
		},
	})
	for i := 0; i < count; i++ {
		key := "key-" + strconv.Itoa(i)
		val, err := rdb.Del(context.Background(), key).Result()
		fmt.Printf("[rdb] key:%s, value:%d, error: %v\n", key, val, err)
	}
	fmt.Println("-----------------------")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		rdb1 := redis.NewRing(&redis.RingOptions{
			Addrs: map[string]string{
				"shard1": "localhost:6379",
				"shard2": "localhost:6380",
			},
		})
		ctx := context.Background()
		for i := 0; i < count; i++ {
			key := "key-" + strconv.Itoa(i)
			value := "value-rdb1-" + strconv.Itoa(i)
			set, err := rdb1.SetNX(ctx, key, value, 0).Result()
			if err != nil {
				fmt.Printf("[rdb1] error: %v", err)
			}
			fmt.Printf("[rdb1] set:%t, key:%s, value:%s\n", set, key, value)
		}
		wg.Done()
	}()

	go func() {
		rdb2 := redis.NewRing(&redis.RingOptions{
			Addrs: map[string]string{
				"shard1": "localhost:6379",
				"shard2": "localhost:6380",
			},
		})
		ctx := context.Background()
		for i := 0; i < count; i++ {
			key := "key-" + strconv.Itoa(i)
			value := "value-rdb2-" + strconv.Itoa(i)
			set, err := rdb2.SetNX(ctx, key, value, 0).Result()
			if err != nil {
				fmt.Printf("[rdb2] error: %v", err)
			}
			fmt.Printf("[rdb2] set:%t, key:%s, value:%s\n", set, key, value)
		}
		wg.Done()
	}()

	// go func() {
	// 	rdb3 := redis.NewRing(&redis.RingOptions{
	// 		Addrs: map[string]string{
	// 			"shard1": "localhost:6379",
	// 			"shard2": "localhost:6380",
	// 		},
	// 	})
	// 	ctx := context.Background()
	// 	for i := 0; i < count; i++ {
	// 		key := "key-" + strconv.Itoa(i)
	// 		value := "value-rdb3-" + strconv.Itoa(i)
	// 		set, err := rdb3.SetNX(ctx, key, value, 0).Result()
	// 		if err != nil {
	// 			fmt.Printf("[rdb3] error: %v", err)
	// 		}
	// 		fmt.Printf("[rdb3] set:%t, key:%s, value:%s\n", set, key, value)
	// 	}
	// 	wg.Done()
	// }()

	wg.Wait()

	fmt.Println("--------- get -----------")
	for i := 0; i < count; i++ {
		key := "key-" + strconv.Itoa(i)
		val, err := rdb.Get(context.Background(), key).Result()
		fmt.Printf("key:%s, value:%s, error: %v\n", key, val, err)
	}
}
