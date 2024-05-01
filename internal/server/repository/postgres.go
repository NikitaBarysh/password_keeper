package repository

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"password_keeper/internal/common/logger"
)

func InitDataBase(ctx context.Context, dsn string, logging *logger.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		logging.Error("err to create DB:\n", err)
		return nil, fmt.Errorf("err to create DB: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		logging.Error("err to ping DB: \n", err)
		return nil, fmt.Errorf("err to ping DB: %w", err)
	}

	return db, nil
}
