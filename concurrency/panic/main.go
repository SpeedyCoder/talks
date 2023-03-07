package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
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

// START1 OMIT
func JobWithCtx(ctx context.Context, jobID int) error {
	select {
	case <-ctx.Done():
		fmt.Printf("context cancelled job %v terminating\n", jobID)
		return ctx.Err()
	case <-time.After(time.Second * time.Duration(rand.Intn(3))):
	}
	if rand.Intn(12) == jobID {
		fmt.Printf("Job %v failed.\n", jobID)
		panic(fmt.Errorf("job %v failed", jobID)) // HL
	}

	fmt.Printf("Job %v done.\n", jobID)
	return nil
}

// END1 OMIT

// START2 OMIT
func main() {
	defer func() { // HL
		if err := recover(); err != nil { // HL
			fmt.Printf("encountered panic: %v\n", err) // HL
		} // HL
	}() // HL
	eg, ctx := errgroup.WithContext(NewCtx())

	for i := 0; i < 10; i++ {
		jobID := i
		eg.Go(func() error {
			return JobWithCtx(ctx, jobID)
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Println("Terminated with error:", err)
	}
}

// END2 OMIT
