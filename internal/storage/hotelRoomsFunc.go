package storage

import (
	"bookings/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

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
