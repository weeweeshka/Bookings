package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upHotelRooms, downHotelRooms)
}

func upHotelRooms(tx *sql.Tx) error {
	const op = "migrations.001_hotel.upHotelRooms"

	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS hotel_rooms(
	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	hotel_id INTEGER NOT NULL,
	rooms INTEGER,
	meals BOOLEAN,
	bar BOOLEAN,
	service BOOLEAN,
	busy BOOLEAN,
	FOREIGN KEY (hotel_id) REFERENCES hotels(id))`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func downHotelRooms(tx *sql.Tx) error {
	const op = "migrations.001_hotel.downHotelRooms"

	_, err := tx.Exec(`DROP TABLE hotel_rooms`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
