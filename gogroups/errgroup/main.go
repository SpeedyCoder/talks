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

func JobWithCtx(ctx context.Context, j int) error {
	select {
	case <-ctx.Done(): // HL
		return nil
	case <-time.After(time.Second * time.Duration(rand.Intn(3))):
	}
	if rand.Intn(12) == j {
		fmt.Printf("Job %v failed.\n", j)
		return fmt.Errorf("job %v failed", j)
	}

	fmt.Printf("Job %v done.\n", j)
	return nil
}

func main() {
	eg, ctx := errgroup.WithContext(NewCtx())

	for i := 0; i < 10; i++ {
		j := i // HL
		eg.Go(func() error {
			return JobWithCtx(ctx, j)
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Println("Terminated with error:", err)
	}
}
