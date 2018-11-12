package errors

import "net/http"

// ErrorType describes the datastore error
type ErrorType string

// BadParameter represents a bad parameter provided by the user
var BadParameter ErrorType = "BAD_PARAM"

// NotFound represents an error in which the record could not be found
var NotFound ErrorType = "NOT_FOUND"

// UnknownError ... something unknown occured
var UnknownError ErrorType = "UNKNOWN"

// New returns a pointer to a new DatastoreError. If the `err` parameter is `nil`, this function will return `nil` so
// errors can be checked against `nil`
func New(typ ErrorType, err error) *DatastoreError {
	if err != nil {
		return &DatastoreError{
			errorType: typ,
			err:       err,
		}
	}
	return nil
}

// DatastoreError implements Error
type DatastoreError struct {
	errorType ErrorType
	err       error
}

// Code attempts to translate the datastore error to a HTTP status code
func (e *DatastoreError) Code() int {
	switch e.errorType {
	case BadParameter:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case UnknownError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Error returns the error as a string (in order to implement the built-in error interface)
func (e *DatastoreError) Error() string {
	return e.err.Error()
}
