package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Context interface {
	Done() <-chan struct{}
	// ...
}

func NewCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sCh := make(chan os.Signal, 1)
		signal.Notify(sCh, syscall.SIGINT, syscall.SIGTERM)
		<-sCh
		fmt.Println("Context cancelled via signal.")
		cancel()
	}()
	return ctx
}

func JobWithCtx(ctx context.Context, jobID int) error {
	select {
	case <-ctx.Done(): // HL
		fmt.Printf("context cancelled job %v terminating\n", jobID)
		return nil
	case <-time.After(time.Second * time.Duration(rand.Intn(3))):
		// Simulate some processing
	}
	if rand.Intn(12) == jobID {
		return fmt.Errorf("job %v failed", jobID)
	}

	fmt.Printf("Job %v done.\n", jobID)
	return nil
}

func main() {
	ctx := NewCtx() // HL
	errchan := make(chan error, 10)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(jobID int) {
			defer wg.Done()
			if err := JobWithCtx(ctx, jobID); err != nil { // HL
				errchan <- err
			}
		}(i)
	}

	wg.Wait()
	close(errchan)
	for err := range errchan {
		fmt.Println("Encountered error:", err)
	}
}
