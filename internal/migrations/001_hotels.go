package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upHotels, downHotels)
}

func upHotels(tx *sql.Tx) error {
	const op = "migrations.001_hotel.upHotels"

	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS hotels(
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			country TEXT NOT NULL,
			city TEXT NOT NULL,
			hotel_name TEXT NOT NULL UNIQUE,
			stars INTEGER NOT NULL CHECK (stars BETWEEN 1 AND 5)
		)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func downHotels(tx *sql.Tx) error {
	const op = "migrations.001_hotel.downHotels"

	_, err := tx.Exec(`DROP TABLE hotels`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
