package errors

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrorNotFound       = errors.New("not found")
	ErrorInvalidRequest = errors.New("invalid request")
	ErrorUnauthorized   = errors.New("unauthorized")
	ErrorForbidden      = errors.New("forbidden")
	ErrorInternalError  = errors.New("internal error")
	ErrorDuplicateKey   = errors.New("duplicate key")
)

func WrapGormError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrorNotFound
	}

	if isUniqueViolation(err) {
		return ErrorDuplicateKey
	}

	return ErrorInternalError
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
