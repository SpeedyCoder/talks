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

func JobWithCtx(ctx context.Context, j int) error {
	select {
	case <-ctx.Done(): // HL
		return nil
	case <-time.After(time.Second * time.Duration(rand.Intn(3))):
	}
	if rand.Intn(12) == j {
		return fmt.Errorf("job %v failed", j)
	}

	fmt.Printf("Job %v done.\n", j)
	return nil
}

func main() {
	ctx := NewCtx()
	errchan := make(chan error, 10) // HL
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			if err := JobWithCtx(ctx, j); err != nil {
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
