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

func JobWithCtx(ctx context.Context, jobID int) error {
	select {
	case <-ctx.Done(): // HL
		fmt.Printf("context cancelled job %v terminating\n", jobID)
		return nil
	case <-time.After(time.Second * time.Duration(rand.Intn(3))):
	}
	if rand.Intn(12) == jobID {
		fmt.Printf("Job %v failed.\n", jobID)
		return fmt.Errorf("job %v failed", jobID)
	}

	fmt.Printf("Job %v done.\n", jobID)
	return nil
}

func main() {
	eg, ctx := errgroup.WithContext(NewCtx())
	eg.SetLimit(2) // HL

	for i := 0; i < 10; i++ {
		jobID := i // HL
		eg.Go(func() error {
			return JobWithCtx(ctx, jobID)
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Println("Terminated with error:", err)
	}
}
