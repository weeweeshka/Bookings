package storage

import (
	"github.com/jackc/pgx"
)

func New() {
	conn := pgx.Connect()
}
