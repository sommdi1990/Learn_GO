package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	b, _ := json.Marshal(person{Name: "Ali", Age: 30})
	fmt.Println(string(b))
}
