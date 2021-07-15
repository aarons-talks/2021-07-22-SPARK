package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	const numWorkers = 10
	ctx, cancel := context.WithCancel(context.Background())
	resCh := make(chan workerResult)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i, resCh)
	}

	time.Sleep(100 * time.Millisecond)
	result := <-resCh
	cancel()
	fmt.Println(
		"worker",
		result.workerNum,
		"finished with result:",
		result.result,
	)
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
	defer wg.Done()
	res := rand.Intn(1000)
	select {
	case result <- workerResult{
		workerNum: workerNum,
		result:    res,
	}:
	case <-ctx.Done():
		fmt.Println("worker", workerNum, "cancelled")
	}
}
