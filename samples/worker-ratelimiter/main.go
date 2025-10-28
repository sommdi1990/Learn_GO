package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(200 * time.Millisecond) // 5 req/sec
	defer ticker.Stop()

	jobs := make(chan int, 10)
	go func() {
		for i := 1; i <= 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	for j := range jobs {
		<-ticker.C
		fmt.Printf("processed job %d at %s\n", j, time.Now().Format(time.RFC3339Nano))
	}
}
