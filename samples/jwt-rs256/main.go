package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func main() {
	// generate ephemeral RSA key pair (demo only)
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey

	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   "user-1",
		"roles": []string{"admin"},
		"exp":   time.Now().Add(5 * time.Minute).Unix(),
	})
	s, _ := tok.SignedString(priv)
	fmt.Println("token:", s)

	t, err := jwt.Parse(s, func(t *jwt.Token) (any, error) { return pub, nil })
	fmt.Println("valid:", err == nil, t != nil)

	// print public key PEM
	b := x509.MarshalPKCS1PublicKey(pub)
	pem.Encode(os.Stdout, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: b})
}
