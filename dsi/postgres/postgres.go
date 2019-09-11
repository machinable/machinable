package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	// db dependency should me transparent to the application
	"github.com/anothrnick/machinable/dsi/models"

	// postgres driver
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

func (d *Database) mapToQuery(filter map[string]interface{}, validFields map[string]bool, filterString *[]string, args *[]interface{}, index *int) error {
	for key, value := range filter {
		if _, ok := validFields[key]; !ok {
			// not a valid field, move on
			continue
		}

		*args = append(*args, value)
		*filterString = append(*filterString, fmt.Sprintf("%s=$%d", key, *index))
		*index++
	}

	return nil
}

func (d *Database) filterToQuery(filter *models.Filters, validFields map[string]bool, filterString *[]string, args *[]interface{}, index *int) error {
	for key, value := range *filter {
		if _, ok := validFields[key]; !ok {
			// not a valid field, move on
			continue
		}
		for op, i := range value {
			var postgresOp string
			switch op {
			case models.GTE:
				postgresOp = ">="
			case models.GT:
				postgresOp = ">"
			case models.LTE:
				postgresOp = "<="
			case models.LT:
				postgresOp = "<"
			case models.EQ:
				postgresOp = "="
			default:
				return errors.New("invalid operator")
			}

			*args = append(*args, i)
			*filterString = append(*filterString, fmt.Sprintf("%s%s$%d", key, postgresOp, *index))
			*index++
		}
	}

	return nil
}
