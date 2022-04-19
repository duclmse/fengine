package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0,
	})
}

func TestClient(t *testing.T) {
	rdb := connect()

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
