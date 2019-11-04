package database

import (
	"database/sql"
	"os"
)

// Open database commection
func Open() (*sql.DB, error) {
	return sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE"))
}
