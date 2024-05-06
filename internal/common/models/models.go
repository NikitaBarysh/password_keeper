// Package models - пакет, в котором лежат кастомные ошибки
package models

import "errors"

var (
	// ErrNotUniqueLogin - ошибка, если занят логин
	ErrNotUniqueLogin = errors.New("login is busy")
)
