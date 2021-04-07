package pgsql

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewConn(dbURI string) (*sql.DB, error) {
	return sql.Open("postgres", dbURI)
}



