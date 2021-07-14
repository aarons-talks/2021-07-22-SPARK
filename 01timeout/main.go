package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	const timeout = 500 * time.Millisecond
	rand.Seed(time.Now().UnixNano())
	timer := time.NewTimer(timeout)

	ch1 := make(chan time.Duration)
	ch2 := make(chan time.Duration)

	go func() {
		ch1 <- maybeLongFunction()
	}()
	go func() {
		ch2 <- maybeLongFunction()
	}()

	fn1Done := false
	fn2Done := false
	for {
		if fn1Done && fn2Done {
			fmt.Println("both functions are done")
			return
		}
		select {
		case slept1 := <-ch1:
			fmt.Println("function 1 woke up after", slept1)
			fn1Done = true
		case slept2 := <-ch2:
			fmt.Println("function 2 woke up after", slept2)
			fn2Done = true
		case <-timer.C:
			fmt.Println("timed out after", timeout)
			os.Exit(1)
		}
	}
}

func maybeLongFunction() time.Duration {
	toSleep := time.Duration(rand.Intn(600)) * time.Millisecond
	time.Sleep(toSleep)
	return toSleep
}
