package main

import (
	"context"
	"fmt"
	"github.com/deliveroo/safe-go"
)

// START OMIT
type Data struct {
	User   *User
	Orders []Order
}

func fetchData(ctx context.Context, userUUID string) (*Data, error) {
	var data Data // HL
	eg, ctx := safe.GroupWithContext(ctx)
	eg.Go(func() (err error) {
		data.User, err = fetchUser(ctx, userUUID) // HL
		return err
	})
	eg.Go(func() (err error) {
		data.Orders, err = fetchOrders(ctx, userUUID) // HL
		return err
	})

	if err := eg.Wait(); err != nil { // HL
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	return &data, nil
}

// END OMIT

type User struct{}

type Order struct{}

func fetchUser(ctx context.Context, userUUID string) (*User, error) {
	return nil, nil
}

func fetchOrders(ctx context.Context, userUUID string) ([]Order, error) {
	return nil, nil
}
