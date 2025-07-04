package storage

import (
	"bookings/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type Postgres struct {
	conn *pgx.Conn
}

func NewPostgresDb(conf *config.Config) (*Postgres, error) {
	const op = "storage.postgres.NewPostgresDB"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, conf.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: connection failed: %w", op, err)
	}

	// make migration

	sqlDB, err := sql.Open("pgx", conf.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: connection for migrations DB failed: %w", op, err)
	}

	goose.SetDialect("postgres")

	if conf.Reload {
		slog.Info("Start reload db")
		err = goose.DownTo(sqlDB, ".", 0)
		if err != nil {
			return nil, fmt.Errorf("%s: Reload failed: %w", op, err)
		}
	}

	slog.Info("Migration started!")
	err = goose.Up(sqlDB, ".")
	if err != nil {
		return nil, fmt.Errorf("%s: Migration failed: %w", op, err)
	}

	// HOTELS TABLE
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS hotels(
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			country TEXT NOT NULL,
			city TEXT NOT NULL,
			hotel_name TEXT NOT NULL UNIQUE,
			stars INTEGER NOT NULL CHECK (stars BETWEEN 1 AND 5)
		)`)
	if err != nil {
		return nil, fmt.Errorf("%s: create table failed: %w", op, err)
	}

	// CreateHotel stmt
	_, err = conn.Prepare(ctx, "create_hotel", `INSERT INTO hotels(country, city, hotel_name, stars) VALUES($1, $2, $3, $4)`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare create_hotel failed: %w", op, err)
	}

	// GetAllHotels stmt

	_, err = conn.Prepare(ctx, "get_all_hotels", `SELECT id, country, city, hotel_name, stars FROM hotels`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare get_all_hotels failed: %w", op, err)
	}

	// GetHotel stmt

	_, err = conn.Prepare(ctx, "get_hotel", `SELECT id, country, city, hotel_name, stars FROM hotels WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare get_hotel failed: %w", op, err)
	}

	// DeleteHotel stmt

	_, err = conn.Prepare(ctx, "delete_hotel", `DELETE FROM hotels WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare delete_hotel failed: %w", op, err)
	}

	//UpdateHotel stmt

	_, err = conn.Prepare(ctx, "update_hotel", `UPDATE hotels SET country = $1, city = $2, hotel_name = $3, stars = $4 WHERE id = $5`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare update_hotel failed: %w", op, err)
	}

	//HOTELROOMS TABLE
	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS hotel_rooms(
	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	hotel_id INTEGER NOT NULL,
	rooms INTEGER,
	meals BOOLEAN,
	bar BOOLEAN,
	service BOOLEAN,
	busy BOOLEAN,
	FOREIGN KEY (hotel_id) REFERENCES hotels(id))`)
	if err != nil {
		return nil, fmt.Errorf("%s: create table failed: %w", op, err)
	}

	// CreateHotelRoom stmt
	_, err = conn.Prepare(ctx, "create_hotel_room", `INSERT INTO hotel_rooms(hotel_id, rooms, meals, bar, service, busy)
	 VALUES($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare create_hotel_room failed: %w", op, err)
	}

	// GetAllHotelRooms stmt
	_, err = conn.Prepare(ctx, "get_all_hotel_rooms", `SELECT id, hotel_id, rooms, meals, bar, service, busy FROM hotel_rooms`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare get_all_hotel_rooms failed: %w", op, err)
	}

	// GetHotelRoom stmt
	_, err = conn.Prepare(ctx, "get_hotel_room", `SELECT * FROM hotel_rooms WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare get_hotel_room failed: %w", op, err)
	}

	// DeleteHotelRoom stmt
	_, err = conn.Prepare(ctx, "delete_hotel_room", `DELETE FROM hotel_rooms WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare delete_hotel_room failed: %w", op, err)
	}

	// UpdateHotelRoom stmt
	_, err = conn.Prepare(ctx, "update_hotel_room", `UPDATE hotel_rooms
	 SET hotel_id = $1, rooms = $2, meals = $3, bar = $4, service = $5 WHERE id = $6`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare update_hotel_room failed: %w", op, err)
	}

	// VISITORS TABLE

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS visitors(
	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	hotel_id INTEGER NOT NULL,
	hotel_room_id INTEGER NOT NULL,
	first_name TEXT,
	last_name TEXT,
	age INTEGER CHECK (age BETWEEN 18 AND 100),
	FOREIGN KEY (hotel_room_id) REFERENCES hotel_rooms(id))`)
	if err != nil {
		return nil, fmt.Errorf("%s: create table failed: %w", op, err)
	}

	// CreateVisitor stmt

	_, err = conn.Prepare(ctx, "create_visitor", `INSERT INTO visitors(hotel_id, hotel_room_id, first_name, last_name, age) VALUES($1, $2, $3, $4, $5)`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare create_vivstor failed: %w", op, err)
	}

	// GetAllVisitors stmt

	_, err = conn.Prepare(ctx, "get_all_visitors", `SELECT id, hotel_id, hotel_room_id, first_name, last_name, age FROM visitors`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare get_all_visitors failed: %w", op, err)
	}

	// GetVisitor stmt

	_, err = conn.Prepare(ctx, "get_visitor", `SELECT id, hotel_id, hotel_room_id, first_name, last_name, age FROM visitors WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare get_visitor failed: %w", op, err)
	}

	// DeleteVisitor stmt

	_, err = conn.Prepare(ctx, "delete_visitor", `DELETE FROM visitors WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare delete_visitor failed: %w", op, err)
	}

	// UpdateVisitor stmt

	_, err = conn.Prepare(ctx, "update_visitor", `UPDATE visitors SET hotel_id = $1, hotel_room_id = $2, first_name = $3, last_name = $4, age = $5`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare update_visitor failed: %w", op, err)
	}

	return &Postgres{conn: conn}, nil
}
