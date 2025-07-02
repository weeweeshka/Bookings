package storage

import (
	"bookings/internal/config"
	"bookings/internal/models"
	"context"
	"database/sql"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, conf.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: connection failed: %w", op, err)
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
	hotel_room_id INTEGER NOT NULL,
	first_name TEXT,
	last_name TEXT,
	age INTEGER CHECK (age BETWEEN 18 AND 100),
	FOREIGN KEY (hotel_room_id) REFERENCES hotel_rooms(id))`)
	if err != nil {
		return nil, fmt.Errorf("%s: create table failed: %w", op, err)
	}

	return &Postgres{conn: conn}, nil
}

// HOTELS TABLE
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

// HOTEL_ROOMS TABLE

func (pos *Postgres) CreateHotelRoom(hotelId int, rooms int, meals bool, bar bool, service bool, busy bool) error {
	const op = "storage.postgres.CreateHotelRoom"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "create_hotel_room", hotelId, rooms, meals, bar, service, busy)
	if err != nil {
		return fmt.Errorf("%s: exec failed: %w", op, err)
	}
	return nil
}

func (pos *Postgres) GetAllHotelRooms() (string, error) {
	const op = "storage.postgres.GetAllHotelRooms"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := pos.conn.Query(ctx, "get_all_hotel_rooms")
	if err != nil {
		return "", fmt.Errorf("%s: query failed: %w", op, err)
	}

	var hotelRooms []models.HotelRoom
	for rows.Next() {
		var hr models.HotelRoom

		if err = rows.Scan(&hr.Id, &hr.HotelId, &hr.Rooms, &hr.Meals, &hr.Bar, &hr.Services, &hr.Busy); err != nil {
			return "", fmt.Errorf("%s: scan failed: %w", op, err)
		}

		hotelRooms = append(hotelRooms, hr)
	}

	jsonHR, err := json.Marshal(hotelRooms)
	if err != nil {
		return "", fmt.Errorf("%s: marshal failed: %w", op, err)
	}

	return string(jsonHR), nil
}

func (pos *Postgres) GetHotelRoom(id int) (string, error) {
	const op = "storage.postgres.GetHotelRoom"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hr models.HotelRoom
	err := pos.conn.QueryRow(ctx, "get_hotel_room", id).Scan(&hr.Id, &hr.HotelId, &hr.Rooms, &hr.Meals, &hr.Bar, &hr.Services, &hr.Busy)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: not found: %w", op, err)
		}
		return "", fmt.Errorf("%s: query failed: %w", op, err)
	}

	jsonHR, _ := json.Marshal(hr)

	return string(jsonHR), nil
}

func (pos *Postgres) DeleteHotelRoom(id int) error {
	const op = "storage.postgres.DeleteHotelRoom"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "delete_hotel_room", id)
	if err != nil {
		return fmt.Errorf("%s: delete failed: %w", op, err)
	}

	return nil
}

func (pos *Postgres) UpdateHotelRoom(id int, hotelId int, rooms int, meals bool, bar bool, service bool) (string, error) {
	const op = "storage.postgres.UpdateHotelRoom"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "update_hotel", hotelId, rooms, meals, bar, service, id)
	if err != nil {
		return "", fmt.Errorf("%s: update failed: %w", op, err)
	}

	return pos.GetHotelRoom(id)
}
