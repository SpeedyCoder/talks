package main

import "context"

type ExampleJob func()

type RealJob func() error

func work(ctx context.Context, in <-chan string, out chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-in:
			if !ok {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- msg:
			}
		}
	}
}
