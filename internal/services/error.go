package services

import "errors"

var (
	ErrDuplicateKey     = errors.New("duplicate key")
	ErrResourceNotFound = errors.New("resource not found")
	ErrUpdateConflict   = errors.New("update conflict")
)
