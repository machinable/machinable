package models

// NewTranslatedError returns a pointer to a new `TranslatedError`
func NewTranslatedError(code int, err error) *TranslatedError {
	return &TranslatedError{
		Code: code,
		Err:  err,
	}
}

// TranslatedError is a database error translated with HTTP Status code
type TranslatedError struct {
	Code int
	Err  error
}

func (t *TranslatedError) Error() string {
	return t.Err.Error()
}
