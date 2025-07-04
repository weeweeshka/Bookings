package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upHotels, downHotels)
}

func upHotels(tx *sql.Tx) error {

}

func downHotels(tx *sql.Tx) error {

}
