package main

import "fmt"

func main() {
	s := make([]string, 3)
	s[0], s[1], s[2] = "a", "b", "c"
	s = append(s, "d")
	fmt.Println(s)
}
