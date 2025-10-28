package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Token bucket ساده با Redis: 5 درخواست در 10 ثانیه برای هر کلید
func allow(ctx context.Context, rdb *redis.Client, key string, limit, windowSec int) (bool, int64, error) {
	rlKey := fmt.Sprintf("rl:%s", key)
	// افزایش شمارنده و تنظیم TTL اگر نبود
	pipe := rdb.TxPipeline()
	incr := pipe.Incr(ctx, rlKey)
	pipe.Expire(ctx, rlKey, time.Duration(windowSec)*time.Second)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, 0, err
	}
	n := incr.Val()
	if n > int64(limit) {
		ttl, _ := rdb.TTL(ctx, rlKey).Result()
		return false, int64(ttl.Seconds()), nil
	}
	return true, 0, nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	for i := 0; i < 7; i++ {
		ok, retry, err := allow(ctx, rdb, "apikey123", 5, 10)
		if err != nil {
			panic(err)
		}
		if !ok {
			fmt.Println("rate limited, retry after (s):", retry)
		} else {
			fmt.Println("allowed", i)
		}
	}
}
