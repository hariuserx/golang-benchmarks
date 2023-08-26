package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	numWorkers = 5000000
)

type Operation int

const (
	Add Operation = iota
	Delete
)

type KeyOperation struct {
	Operation Operation
	Key       *int
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(0))
	keySet := make(map[*int]bool)
	channel := make(chan KeyOperation)
	done := make(chan bool)

	// Measure time for map operations with Mutex and Channel
	fmt.Println("Map Operations with Unbuffered Channel:")
	startTime := time.Now()
	var wg sync.WaitGroup

	// Start worker goroutines to send keys to the channel
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		k := i
		go func() {
			defer wg.Done()
			channel <- KeyOperation{Operation: Add, Key: &k}
			time.Sleep(100 * time.Millisecond)
			channel <- KeyOperation{Operation: Delete, Key: &k}
		}()
	}

	// Start the key polling goroutine
	go func() {
		defer func() {
			// Sending value to channel
			done <- true
		}()
		l := 0
		for {
			select {
			case key, ok := <-channel:
				if !ok {
					fmt.Printf("Channel closed. Operations done : %v\n", l)
					return
				}
				if key.Operation == Add {
					keySet[key.Key] = true
				} else {
					delete(keySet, key.Key)
				}
				l++
			}
		}
	}()

	// Wait for all worker goroutines to finish
	wg.Wait()
	close(channel)
	<-done

	fmt.Println("Time taken:", time.Since(startTime))
}
