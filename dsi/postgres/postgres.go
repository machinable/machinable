package postgres

import (
	"database/sql"
	"fmt"

	// db dependency should me transparent to the application
	_ "github.com/lib/pq"
)

// Database is a wrapper for the PostgreSQL connection
type Database struct {
	db *sql.DB
}

// New creates and returns a pointer to a new instance of `Database`
func New(user, password, host, database string) (*Database, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}
