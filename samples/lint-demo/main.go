package main

import "fmt"

// bad is intentionally unused to trigger a lint warning (deadcode)
func bad() string {
	return "unused"
}

func main() {
	fmt.Println("lint demo: fix issues then run make lint")
}
