package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func Job(j int) error {
	if rand.Intn(12) == j {
		return fmt.Errorf("job %v failed", j)
	}

	fmt.Printf("Job %v done.\n", j)
	return nil
}

func main() {
	errchan := make(chan error, 10) // HL
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			if err := Job(j); err != nil {
				errchan <- err // HL
			}
		}(i)
	}

	wg.Wait()
	close(errchan)
	for err := range errchan {
		fmt.Println("Encountered error:", err)
	}
}
