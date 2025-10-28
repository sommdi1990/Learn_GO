package main

import (
	"errors"
	"fmt"
)

func mightFail(b bool) error {
	if b {
		return errors.New("failed")
	}
	return nil
}

func main() { fmt.Println(mightFail(true)) }
