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

var mutex sync.Mutex
var keySet map[*int]bool

func addToMap(key *int) {
	mutex.Lock()
	defer mutex.Unlock()

	keySet[key] = true
}

func deleteFromMap(key *int) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(keySet, key)
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(0))
	keySet = make(map[*int]bool)

	// Measure time for map operations with Mutex and Channel
	fmt.Println("Map Operations with Mutex :")
	startTime := time.Now()
	var wg sync.WaitGroup

	// Start worker goroutines to send keys to the channel
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		k := i
		go func() {
			defer wg.Done()
			addToMap(&k)
			time.Sleep(100 * time.Millisecond)
			deleteFromMap(&k)
		}()
	}

	// Wait for all worker goroutines to finish
	wg.Wait()

	fmt.Println("Time taken:", time.Since(startTime))
}
