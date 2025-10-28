package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Println("connected to Redis")

	// SET/GET with TTL
	if err := rdb.Set(ctx, "demo:key", "value", 30*time.Second).Err(); err != nil {
		log.Fatalf("set: %v", err)
	}
	val, err := rdb.Get(ctx, "demo:key").Result()
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Println("demo:key=", val)

	// INCR
	n, err := rdb.Incr(ctx, "demo:counter").Result()
	if err != nil {
		log.Fatalf("incr: %v", err)
	}
	fmt.Println("demo:counter=", n)

	// PUB/SUB
	go func() {
		time.Sleep(200 * time.Millisecond)
		_ = rdb.Publish(context.Background(), "demo:topic", "hello").Err()
	}()
	sub := rdb.Subscribe(ctx, "demo:topic")
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		log.Fatalf("subscribe: %v", err)
	}
	fmt.Printf("received pubsub: %s=%s\n", msg.Channel, msg.Payload)
}
