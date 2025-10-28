package main

import (
	"context"
	"fmt"
	"log"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	ctx := context.Background()
	conn, err := ch.Open(&ch.Options{
		Addr: []string{"localhost:9000"},
		Auth: ch.Auth{
			Database: "demo",
			Username: "demo",
			Password: "demo",
		},
		DialTimeout: 5 * time.Second,
		Compression: &ch.Compression{Method: ch.CompressionLZ4},
		Settings: ch.Settings{
			"allow_experimental_object_type": 1,
		},
	})
	if err != nil {
		log.Fatalf("connect: %v", err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}
	log.Println("connected to ClickHouse")

	ddl := `
    CREATE TABLE IF NOT EXISTS events (
        event_date Date DEFAULT today(),
        user_id    UInt64,
        event_name LowCardinality(String),
        props      String
    ) ENGINE = MergeTree
    PARTITION BY toYYYYMM(event_date)
    ORDER BY (event_date, user_id)
    `
	if err := conn.Exec(ctx, ddl); err != nil {
		log.Fatalf("create table: %v", err)
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO events (user_id, event_name, props)")
	if err != nil {
		log.Fatalf("batch: %v", err)
	}
	for i := 1; i <= 5; i++ {
		if err := batch.Append(uint64(i), "login", fmt.Sprintf("{\"n\":%d}", i)); err != nil {
			log.Fatalf("append: %v", err)
		}
	}
	if err := batch.Send(); err != nil {
		log.Fatalf("send: %v", err)
	}
	log.Println("inserted 5 rows")

	type row struct {
		UserID    uint64
		EventName string
		Props     string
	}

	rows, err := conn.Query(ctx, "SELECT user_id, event_name, props FROM events WHERE event_name = 'login' ORDER BY user_id LIMIT 5")
	if err != nil {
		log.Fatalf("query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var r row
		if err := rows.ScanStruct(&r); err != nil {
			log.Fatalf("scan: %v", err)
		}
		fmt.Printf("user=%d event=%s props=%s\n", r.UserID, r.EventName, r.Props)
	}
}
