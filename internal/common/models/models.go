package models

import "errors"

var (
	ErrNotUniqueLogin = errors.New("login is busy")
)
