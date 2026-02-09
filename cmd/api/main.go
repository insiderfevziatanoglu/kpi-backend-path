package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fevziatanoglu/test-go-project/internal/config"
	"github.com/fevziatanoglu/test-go-project/internal/models" 
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println("starting...")
	fmt.Printf("env: %s\n", cfg.AppEnv)
	fmt.Printf("port: %s\n", cfg.ServerPort)

	if cfg.DBUrl != "" {
		log.Println("connected.")
	}

	fmt.Println("\n test json")

	testUser := models.User{
		ID:           101,
		Username:     "fevzi_test",
		Email:        "fevzi@example.com",
		PasswordHash: "testpassword",
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	jsonData, err := json.MarshalIndent(testUser, "", "  ")
	if err != nil {
		log.Fatal("JSON  error:", err)
	}

	fmt.Println(string(jsonData))
	fmt.Println("---------------------------------")

}
