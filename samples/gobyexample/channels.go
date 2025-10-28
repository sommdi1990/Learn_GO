package main

import "fmt"

func main() {
	ch := make(chan string)
	go func() { ch <- "ping" }()
	fmt.Println(<-ch)
}
