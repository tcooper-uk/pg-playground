package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connections struct {
	Primary *pgxpool.Pool
	Replica *pgxpool.Pool
}

func Connect(ctx context.Context, primaryDSN, replicaDSN string) (*Connections, error) {
	primary, err := pgxpool.New(ctx, primaryDSN)
	if err != nil {
		return nil, fmt.Errorf("primary: %w", err)
	}
	if err := primary.Ping(ctx); err != nil {
		return nil, fmt.Errorf("primary ping: %w", err)
	}

	replica, err := pgxpool.New(ctx, replicaDSN)
	if err != nil {
		return nil, fmt.Errorf("replica: %w", err)
	}
	if err := replica.Ping(ctx); err != nil {
		return nil, fmt.Errorf("replica ping: %w", err)
	}

	return &Connections{Primary: primary, Replica: replica}, nil
}

func (c *Connections) Close() {
	c.Primary.Close()
	c.Replica.Close()
}
