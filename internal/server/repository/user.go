package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"password_keeper/internal/common/entity"
)

type AuthRepository struct {
	rep *sqlx.DB
}

func NewAuthRepository(newDB *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		rep: newDB,
	}
}

func (r *AuthRepository) SetUserDB(ctx context.Context, user entity.User) (int, error) {
	var id int

	tx, err := r.rep.Beginx()
	if err != nil {
		return 0, fmt.Errorf("err to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO users (login, password) VALUES ($1,$2) RETURNING id`,
		user.Login, user.Password)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to do exec int DB: %w", err)
	}

	row := tx.QueryRowxContext(ctx, "SELECT id FROM  users WHERE login=$1", user.Login)
	if row.Err() != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to get id: %w", row.Err())
	}

	if err = row.Scan(&id); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
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
		return 0, fmt.Errorf("err to begin transaction: %w", err)
	}

	row := tx.QueryRowxContext(ctx, `SELECT id FROM users WHERE login=$1 AND password=$2`, user.Login, user.Password)
	if row.Err() != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to get id: %w", row.Err())
	}

	if err = row.Scan(&id); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return 0, fmt.Errorf("err to do Rollback: %w", errRollback)
		}
		return 0, fmt.Errorf("err to scan id: %w", err)
	}

	return id, tx.Commit()

}

func (r *AuthRepository) Validate(ctx context.Context, username string) error {
	var id int

	row := r.rep.QueryRowxContext(ctx, "SELECT id FROM users WHERE login=$1", username)
	if row.Err() != nil {
		return fmt.Errorf("err to get id: %w", row.Err())
	}

	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("err to scan id: %w", err)
	}

	return fmt.Errorf("err to get id: %w", row.Err())
}
