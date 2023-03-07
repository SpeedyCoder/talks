package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func Job(jobID int) error {
	if rand.Intn(12) == jobID {
		return fmt.Errorf("job %v failed", jobID)
	}

	fmt.Printf("Job %v done.\n", jobID)
	return nil
}

func main() {
	errchan := make(chan error, 10) // HL
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(jobID int) {
			defer wg.Done()
			if err := Job(jobID); err != nil {
				errchan <- err // HL
			}
		}(i)
	}

	wg.Wait()
	close(errchan)             // HL
	for err := range errchan { // HL
		fmt.Println("Encountered error:", err)
	}
}
