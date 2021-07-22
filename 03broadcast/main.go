package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	bcastCh := make(chan struct{})
	const numWaiters = 10
	for i := 0; i < numWaiters; i++ {
		wg.Add(1)
		go waiter(i, bcastCh, &wg)
	}

	fmt.Println("going to sleep...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("broadcasting to all", numWaiters, "waiters")
	close(bcastCh)
	wg.Wait()
	fmt.Println("done")
}

func waiter(funcNum int, ch <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ch
	fmt.Println("function", funcNum, "is woken up!")
}
