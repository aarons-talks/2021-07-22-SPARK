package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	const numWorkers = 10
	const timeoutDur = 1 * time.Second
	fmt.Println("creating context with timeout", timeoutDur)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDur)
	defer cancel()
	resCh := make(chan workerResult)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i, resCh)
	}

	for i := 0; i < numWorkers; i++ {
		sleepDur := time.Duration(rand.Intn(500)) * time.Millisecond
		time.Sleep(sleepDur)
		select {
		case result := <-resCh:
			fmt.Println(
				"worker",
				result.workerNum,
				"finished with result:",
				result.result,
			)
		case <-ctx.Done():
			fmt.Println("context done before all workers finished")
			os.Exit(1)
		}
	}
	wg.Wait()
}

type workerResult struct {
	workerNum int
	result    int
}

func worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	workerNum int,
	result chan<- workerResult,
) {
	start := time.Now()
	defer wg.Done()
	res := rand.Intn(1000)
	select {
	case result <- workerResult{
		workerNum: workerNum,
		result:    res,
	}:
	case <-ctx.Done():
		fmt.Println("worker", workerNum, "cancelled after", time.Since(start))
	}
}
