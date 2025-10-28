package main

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secret = []byte("dev-secret")

func issue() (string, error) {
	claims := jwt.MapClaims{
		"sub":   "user-1",
		"roles": []string{"reader"},
		"exp":   time.Now().Add(10 * time.Minute).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func verify(tok string) (*jwt.Token, error) {
	return jwt.Parse(tok, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("alg")
		}
		return secret, nil
	})
}

func main() {
	tok, _ := issue()
	fmt.Println("token:", tok)
	t, err := verify(tok)
	fmt.Println("valid:", err == nil, t != nil)
}
