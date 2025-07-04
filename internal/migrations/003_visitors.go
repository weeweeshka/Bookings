package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upVisitors, downVisitors)
}

func upVisitors(tx *sql.Tx) error {
	const op = "migrations.001_hotel.upVisitors"

	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS visitors(
	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	hotel_id INTEGER NOT NULL,
	hotel_room_id INTEGER NOT NULL,
	first_name TEXT,
	last_name TEXT,
	age INTEGER CHECK (age BETWEEN 18 AND 100),
	FOREIGN KEY (hotel_room_id) REFERENCES hotel_rooms(id))`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func downVisitors(tx *sql.Tx) error {
	const op = "migrations.001_hotel.downVisitors"

	_, err := tx.Exec(`DROP TABLE visitors`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
