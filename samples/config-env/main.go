package main

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	HTTPAddr string
	Env      string
}

func loadConfig() Config {
	cfg := Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		Env:      getEnv("APP_ENV", "dev"),
	}
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	cfg := loadConfig()
	log.Printf("config: addr=%s env=%s\n", cfg.HTTPAddr, cfg.Env)
	fmt.Println("OK")
}
