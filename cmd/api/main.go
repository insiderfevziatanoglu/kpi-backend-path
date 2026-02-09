package main

import (
	"fmt"
	"log"

	"github.com/fevziatanoglu/test-go-project/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println("starting...")
	fmt.Printf("env: %s\n", cfg.AppEnv)
	fmt.Printf("port: %s\n", cfg.ServerPort)

	if cfg.DBUrl != "" {
		log.Println("connected.")
	}
}