package storage

import (
	"bookings/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func (pos *Postgres) CreateHotel(country string, city string, hotelName string, stars int) error {
	const op = "storage.postgres.CreateHotel"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "create_hotel", country, city, hotelName, stars)
	if err != nil {
		return fmt.Errorf("%s: exec failed: %w", op, err)
	}

	return nil
}

func (pos *Postgres) GetAllHotels() (string, error) {
	const op = "storage.postgres.GetAllHotels"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := pos.conn.Query(ctx, "get_all_hotels")
	if err != nil {
		return "", fmt.Errorf("%s: scan failed: %w", op, err)
	}

	var hotels []models.Hotel
	for rows.Next() {
		var h models.Hotel
		if err := rows.Scan(&h.Country, &h.City, &h.HotelName, &h.Stars); err != nil {
			return "", fmt.Errorf("%s: scan failed: %w", op, err)
		}
		hotels = append(hotels, h)
	}

	jsonData, err := json.Marshal(hotels)
	if err != nil {
		return "", fmt.Errorf("%s: marshal failed: %w", op, err)
	}

	return string(jsonData), nil
}

func (pos *Postgres) GetHotel(id int) (string, error) {
	const op = "storage.postgres.GetHotel"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hotel models.Hotel
	err := pos.conn.QueryRow(ctx, "get_hotel", id).Scan(&hotel.Id, &hotel.Country, &hotel.City, &hotel.HotelName, &hotel.Stars)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: not found: %w", op, err)
		}
		return "", fmt.Errorf("%s: query failed: %w", op, err)
	}

	jsonHotel, _ := json.Marshal(hotel)

	return string(jsonHotel), nil
}

func (pos *Postgres) DeleteHotel(id int) error {
	const op = "storage.postgres.DeleteHotel"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "delete_hotel", id)
	if err != nil {
		return fmt.Errorf("%s: delete failed: %w", op, err)
	}

	return nil
}

func (pos *Postgres) UpdateHotel(id int, country string, city string, hotelName string, stars int) (string, error) {
	const op = "storage.postgres.UpdateHotel"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "update_hotel", country, city, hotelName, stars, id)
	if err != nil {
		return "", fmt.Errorf("%s: update failed: %w", op, err)
	}

	return pos.GetHotel(id)
}
