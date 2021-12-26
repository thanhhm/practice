package resources

import (
	"database/sql"

	_ "github.com/godror/godror"
)

func NewConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}
