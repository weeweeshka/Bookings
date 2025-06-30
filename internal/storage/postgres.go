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

	return &Postgres{conn: conn}, nil
}
