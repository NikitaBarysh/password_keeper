package models

import "errors"

var (
	ErrNotUniqueLogin = errors.New("login is busy")
	ErrScan           = errors.New("err: err to scan id: sql: no rows in result set ")
)
