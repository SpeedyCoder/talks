package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup // HL
	for i := 0; i < 10; i++ {
		wg.Add(1) // HL
		go func() {
			defer wg.Done() // HL
			fmt.Printf("This is job: %v\n", i)
		}()
	}
	wg.Wait() // HL
}
