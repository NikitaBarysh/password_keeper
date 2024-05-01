package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/common/logger"
)

//go:generate mockgen -source ${GOFILE} -destination mock.go -package ${GOPACKAGE}

type AuthorizationRepository interface {
	SetUserDB(ctx context.Context, user entity.User) (int, error)
	GetUserFromDB(ctx context.Context, user entity.User) (int, error)
}

type DataRepositoryInterface interface {
	SetRepData(ctx context.Context, id int, text []byte, eventType string) error
	GetRepData(ctx context.Context, id int, eventType string) ([]byte, error)
	DeleteRepData(ctx context.Context, id int, eventType string) error
}

type Repository struct {
	AuthorizationRepository
	DataRepositoryInterface
}

func NewRepository(db *sqlx.DB, log *logger.Logger) *Repository {
	return &Repository{
		AuthorizationRepository: NewAuthRepository(db, log),
		DataRepositoryInterface: NewDataRepository(db, log),
	}
}
