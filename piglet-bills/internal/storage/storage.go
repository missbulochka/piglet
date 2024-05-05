package storage

import "errors"

var (
	ErrBillExists   = errors.New("Bill already exists")
	ErrBillNotFound = errors.New("Bill not found")
)
