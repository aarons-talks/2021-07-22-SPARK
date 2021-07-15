package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	grp, _ := errgroup.WithContext(ctx)

	const numToRun = 10
	for i := 0; i < numToRun; i++ {
		funcNum := i
		grp.Go(func() error {
			sleptDur, err := maybeLongFunction(funcNum)
			if err != nil {
				return err
			}
			fmt.Println("function", funcNum, "slept for", sleptDur)
			return nil
		})
	}
	result := grp.Wait()
	if result == nil {
		fmt.Println("success!")
	} else {
		fmt.Println("error:", result)
		os.Exit(1)
	}
}

func maybeLongFunction(funcNum int) (time.Duration, error) {
	toSleep := time.Duration(rand.Intn(600)) * time.Millisecond
	time.Sleep(toSleep)
	var err error
	const threshold = 400 * time.Millisecond
	if toSleep > threshold {
		err = fmt.Errorf(
			"function %d slept for more than %v",
			funcNum,
			threshold,
		)
	}
	return toSleep, err
}
