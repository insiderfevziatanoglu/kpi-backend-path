package processor

import (
	"log"
	"sync"
	"time"
)


type BatchHandler func([]int64) error

type BatchProcessor struct {
	InputChan   chan int64       
	BatchSize   int               
	Timeout     time.Duration   
	Handler     BatchHandler      
	quit        chan bool         
	wg          sync.WaitGroup
}

func NewBatchProcessor(size int, timeout time.Duration, handler BatchHandler) *BatchProcessor {
	return &BatchProcessor{
		InputChan: make(chan int64, size*10),
		BatchSize: size,
		Timeout:   timeout,
		Handler:   handler,
		quit:      make(chan bool),
	}
}

func (bp *BatchProcessor) Start() {
	bp.wg.Add(1)
	go bp.loop() 
	log.Println("batch processor started.")
}

func (bp *BatchProcessor) loop() {
	defer bp.wg.Done()

	var batch []int64
	ticker := time.NewTicker(bp.Timeout)

	for {
		select {
		case item := <-bp.InputChan:
			batch = append(batch, item)
			
			if len(batch) >= bp.BatchSize {
				log.Printf("batch full %d", len(batch))
				bp.processBatch(batch)
				batch = nil
				ticker.Reset(bp.Timeout)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				log.Printf("timeout reached %d", len(batch))
				bp.processBatch(batch)
				batch = nil
			}

		case <-bp.quit:
			if len(batch) > 0 {
				log.Printf("shutting down %d", len(batch))
				bp.processBatch(batch)
			}
			return
		}
	}
}

func (bp *BatchProcessor) processBatch(batch []int64) {
	// update/insert
	if err := bp.Handler(batch); err != nil {
		log.Printf("batch error %v", err)
	}
}

func (bp *BatchProcessor) Add(id int64) {
	bp.InputChan <- id
}

func (bp *BatchProcessor) Stop() {
	bp.quit <- true
	bp.wg.Wait()
	log.Println("batch processor stopped.")
}