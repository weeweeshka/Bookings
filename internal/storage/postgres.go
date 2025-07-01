package storage

import (
	"bookings/internal/config"
	"bookings/internal/models"
	"context"
	"encoding/json"
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

	_, err = conn.Prepare(ctx, "hotels_table",
		`CREATE TABLE IF NOT EXISTS hotels(
	id INTEGER PRIMARY KEY,
	country TEXT NOT NULL,
	city TEXT NOT NULL,
	hotel_name TEXT NOT NULL,
	stars INTEGER NOT NULL CHECK (stars BETWEEN 1 AND 5))`)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to create prepare: %w", op, err)
	}

	_, err = conn.Exec(ctx, "hotels_table")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create hotels table: %w", op, err)
	}

	return &Postgres{conn: conn}, nil
}

func (pos *Postgres) CreateHotel(country string, city string, hotelName string, stars int) error {
	const op = "storage.postgres.CreateHotel"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Prepare(ctx, "create_hotel",
		`INSERT INTO hotels(country, city, hotel_name, stars)
	VALUES($1,$2,$3,$4)`)
	if err != nil {
		return fmt.Errorf("%s: failed in prepare: %w", op, err)
	}

	_, err = pos.conn.Exec(ctx, "create_hotel", country, city, hotelName, stars)
	if err != nil {
		return fmt.Errorf("%s: failed to insert data into hotels_table: %w", op, err)
	}

	return nil
}

func (pos *Postgres) GetAllHotels(id int) (string, error) {
	const op = "storage.postgres.GetAllHotels"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := pos.conn.Query(ctx, `SELECT id, country, city, hotel_name, stars FROM hotels WHERE id = ?`)
	if err != nil {
		return "", fmt.Errorf("%s: failed in query: %w", op, err)
	}

	var hotels []models.Hotel
	for rows.Next() {
		var h models.Hotel

		if err := rows.Scan(&h.Country, &h.City, &h.HotelName, &h.Stars); err != nil {
			return "", fmt.Errorf("%s: failed in query: %w", op, err)
		}

		hotels = append(hotels, h)
	}

	jsonHotel, err := json.Marshal(hotels)
	if err != nil {
		return "", fmt.Errorf("%s: failed in marshal: %w", op, err)
	}

	return string(jsonHotel), nil
}

func (pos *Postgres) GetHotel(id int) (models.Hotel, error) {
	const op = "storage.postgres.GetHotel"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Prepare(ctx, "get_hotel",
		`SELECT id, country, city, hotel_name, stars FROM hotels WHERE id = ?`)
	if err != nil {
		return models.Hotel{}, fmt.Errorf("%s: failed in prepare: %w", op, err)
	}

	var hotel models.Hotel
	_ = pos.conn.QueryRow(ctx, "get_hotel", id).Scan(&hotel.Country, &hotel.City, &hotel.HotelName, &hotel.Stars)

	return hotel, nil

}

func (pos *Postgres) DeleteHotel(id int) error {
	const op = "storage.postgres.DeleteHotel"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Prepare(ctx, "delete_hotel",
		`DELETE FROM hotels WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("%s: failed in prepare: %w", op, err)
	}

	_, err = pos.conn.Exec(ctx, "delete_hotel", id)
	if err != nil {
		return fmt.Errorf("%s: failed in exec: %w", op, err)
	}

	return nil
}

func (pos *Postgres) UpdateHotel(id int, country string, city string, hotelName string, stars int) (models.Hotel, error) {
	const op = "storage.postgres.UpdateHotel"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Prepare(ctx, "update_hotel",
		`UPDATE hotels SET country = ?, city = ?, hotel_name = ?, stars = ? WHERE id = ?`)
	if err != nil {
		return models.Hotel{}, fmt.Errorf("%s: failed in prepare: %w", op, err)
	}

	_, err = pos.conn.Exec(ctx, "update_hotel", country, city, hotelName, stars, id)
	if err != nil {
		return models.Hotel{}, fmt.Errorf("%s: failed in prepare: %w", op, err)
	}

	getHotel, err := pos.GetHotel(id)
	if err != nil {
		return models.Hotel{}, fmt.Errorf("%s: failed in get: %w", op, err)
	}

	return getHotel, nil

}
