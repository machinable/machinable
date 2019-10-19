package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/anothrNick/machinable/dsi/models"
	"github.com/lib/pq"
)

// TranslateError attempts to translate the database specific error to a simple `error` to return to the user.
func (p *Database) TranslateError(err error) *models.TranslatedError {
	originalError := err.Error()

	// log original error

	if err, ok := err.(*pq.Error); ok {
		// postgres specific errors
		switch err.Code {
		case "23505":
			return models.NewTranslatedError(http.StatusBadRequest, errors.New("key already exists"))
		case "22023":
			return models.NewTranslatedError(http.StatusBadRequest, errors.New("key already exists"))
		default:
			fmt.Println(err)
			fmt.Println(err.Code)
			return models.NewTranslatedError(http.StatusInternalServerError, errors.New(strings.TrimPrefix(err.Error(), "pq: ")))
		}
	} else {
		// generic sql package errors
		switch originalError {
		case sql.ErrConnDone.Error():
			return models.NewTranslatedError(http.StatusInternalServerError, errors.New("internal server error"))
		case sql.ErrNoRows.Error():
			return models.NewTranslatedError(http.StatusNotFound, errors.New("not found"))
		case sql.ErrTxDone.Error():
			return models.NewTranslatedError(http.StatusInternalServerError, errors.New("internal server error"))
		}
	}

	return models.NewTranslatedError(http.StatusInternalServerError, errors.New("internal server error"))
}
