package processor

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fevziatanoglu/test-go-project/internal/models"
)

type TransactionJob struct {
	Tx *models.Transaction
}

type WorkerPool struct {
	JobQueue     chan TransactionJob
	WorkerCount  int
	wg           sync.WaitGroup
	quit         chan bool
	ProcessedCount uint64 
	FailedCount    uint64
}

func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	return &WorkerPool{
		JobQueue:    make(chan TransactionJob, queueSize),
		WorkerCount: workerCount,
		quit:        make(chan bool),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.WorkerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i + 1)
	}
	log.Printf("Worker pool started. Workers ready: %d", wp.WorkerCount)
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case job := <-wp.JobQueue:
			log.Printf("Worker %d: Processing amount %.2f", id, job.Tx.Amount)
			
			time.Sleep(500 * time.Millisecond)
			
			atomic.AddUint64(&wp.ProcessedCount, 1)
			
			log.Printf("Worker %d: Done", id)

		case <-wp.quit:
			log.Printf("Worker %d: Stopping", id)
			return
		}
	}
}

func (wp *WorkerPool) Submit(tx *models.Transaction) {
	job := TransactionJob{Tx: tx}
	wp.JobQueue <- job
}

func (wp *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")
	
	for i := 0; i < wp.WorkerCount; i++ {
		wp.quit <- true
	}
	
	wp.wg.Wait()
	log.Println("Worker pool stopped.")
}

func (wp *WorkerPool) GetStats() (uint64, uint64) {
	processed := atomic.LoadUint64(&wp.ProcessedCount)
	failed := atomic.LoadUint64(&wp.FailedCount)
	
	return processed, failed
}