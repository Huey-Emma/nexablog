package repository

import (
	"errors"
	"strings"
)

var (
	ErrDuplicateKey     = errors.New("duplicate key")
	ErrResourceNotFound = errors.New("resource not found")
	ErrUpdateConflict   = errors.New("update conflict")
)

func DuplicateKey(e error) bool {
	return strings.Contains(e.Error(), "duplicate")
}
