package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/common/logger"
)

type AuthRepository struct {
	rep     *sqlx.DB
	logging *logger.Logger
}

func NewAuthRepository(newDB *sqlx.DB, log *logger.Logger) *AuthRepository {
	return &AuthRepository{
		rep:     newDB,
		logging: log,
	}
}

func (r *AuthRepository) SetUserDB(ctx context.Context, user entity.User) (int, error) {
	var id int

	tx, err := r.rep.Beginx()
	if err != nil {
		r.logging.Error("err to begin transaction: ", err)
		return 0, fmt.Errorf("err to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO users (login, password) VALUES ($1,$2) RETURNING id`,
		user.Login, user.Password)
	if err != nil {
		r.logging.Error("err to exec login and password into DB: ", err)
		if errRollback := tx.Rollback(); errRollback != nil {
			r.logging.Error("err to do rollback after err to exec into DB: ", errRollback)
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to do exec int DB: %w", err)
	}

	row := tx.QueryRowxContext(ctx, "SELECT id FROM  users WHERE login=$1", user.Login)
	if row.Err() != nil {
		r.logging.Error("err to get id: ", row.Err())
		if errRollback := tx.Rollback(); errRollback != nil {
			r.logging.Error("err to do rollback after err to get id: ", errRollback)
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to get id: %w", row.Err())
	}

	if err = row.Scan(&id); err != nil {
		r.logging.Error("err to scan id: ", err)
		if errRollback := tx.Rollback(); errRollback != nil {
			r.logging.Error("err to do rollback after err to get id: ", errRollback)
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to scan id: %w", err)
	}

	return id, tx.Commit()
}

func (r *AuthRepository) GetUserFromDB(ctx context.Context, user entity.User) (int, error) {
	var id int

	tx, err := r.rep.Beginx()
	if err != nil {
		r.logging.Error("err to begin transaction: ", err)
		return 0, fmt.Errorf("err to begin transaction: %w", err)
	}

	row := tx.QueryRowxContext(ctx, `SELECT id FROM users WHERE login=$1 AND password=$2`, user.Login, user.Password)
	if row.Err() != nil {
		r.logging.Error("err to get id: ", row.Err())
		if errRollback := tx.Rollback(); errRollback != nil {
			r.logging.Error("err to do rollback after err to get id: ", errRollback)
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to get id: %w", row.Err())
	}

	if err = row.Scan(&id); err != nil {
		r.logging.Error("err to scan id: ", err)
		if errRollback := tx.Rollback(); errRollback != nil {
			r.logging.Error("err to do rollback after err to get id: ", errRollback)
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to scan id: %w", err)
	}

	return id, tx.Commit()

}
