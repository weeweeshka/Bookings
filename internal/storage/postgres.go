package storage

import (
	"bookings/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	conn *pgx.Conn
}

func NewPostgresDb(conf *config.Config) (*Postgres, error) {
	const op = "storage.postgres.NewPostgresDB"

	connConfig, err := pgx.ParseConfig(conf.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s cannot parse connConfig : %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("%s cannot connect : %w", op, err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}

	_, err = conn.Prepare(ctx, "create_bookings_table",
		`CREATE TABLE IF NOT EXISTS bookings(
	id INTEGER PRIMARY KEY,
	country TEXT NOT NULL,
	city TEXT NOT NULL,
	hotel_name TEXT NOT NULL,
	stars INTEGER NOT NULL CHECK (stars BETWEEN 1 AND 5))`)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to create prepare: %w", op, err)
	}

	_, err = conn.Exec(ctx, "create_bookings_table")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create bookings table: %w", op, err)
	}

	return &Postgres{conn: conn}, nil
}

// func (pos *Postgres) CreateBooking() error {

// }
