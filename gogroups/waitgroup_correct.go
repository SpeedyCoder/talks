package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(jobID int) { // HL
			defer wg.Done()
			fmt.Printf("This is job: %v\n", jobID) // HL
		}(i) // HL
	}
	wg.Wait()
}
