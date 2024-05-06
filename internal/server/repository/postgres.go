// Package repository - пакет в котором сохраняем данные в базу данных
package repository

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// InitDataBase - подключаемся к базе данных
func InitDataBase(ctx context.Context, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("InitDataBase: err to create DB: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("InitDataBase: err to ping DB: %w", err)
	}

	return db, nil
}
