package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fevziatanoglu/test-go-project/internal/config"
	"github.com/fevziatanoglu/test-go-project/internal/models"
	"github.com/fevziatanoglu/test-go-project/internal/processor"
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

	/*
		fmt.Println("--- Worker Pool Test ---")

		pool := processor.NewWorkerPool(3, 5)

		pool.Start()

		for i := 1; i <= 6; i++ {
			tx := &models.Transaction{
				ID:     int64(i),
				Amount: float64(i * 100),
			}
			fmt.Printf("Transaction %d submitted\n", i)
			pool.Submit(tx)
		}

		time.Sleep(4 * time.Second)

		pool.Stop()
	*/

	fmt.Println("--- Atomic Counter Test ---")

	pool := processor.NewWorkerPool(3, 10)
	pool.Start()

	for i := 0; i < 10; i++ {
		tx := &models.Transaction{ID: int64(i), Amount: 10.0}
		pool.Submit(tx)
	}

	time.Sleep(2 * time.Second)

	success, failed := pool.GetStats()

	fmt.Println("------------------------------------------------")
	fmt.Printf("Total Processed: %d\n", success)
	fmt.Printf("Total Failed:    %d\n", failed)
	fmt.Println("------------------------------------------------")

	pool.Stop()

	fmt.Println("--- batch test ---")

	myHandler := func(ids []int64) error {
		fmt.Printf("send data: %v\n", ids)
		time.Sleep(500 * time.Millisecond)
		return nil
	}

	batchProc := processor.NewBatchProcessor(5, 3*time.Second, myHandler)

	batchProc.Start()


	for i := 1; i <= 12; i++ {
		fmt.Printf("item added: %d\n", i)
		batchProc.Add(int64(i))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("timeout")
	time.Sleep(4 * time.Second)


	batchProc.Stop()
}
