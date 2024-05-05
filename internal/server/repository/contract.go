package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"password_keeper/internal/common/entity"
)

//go:generate mockgen -source ${GOFILE} -destination mock.go -package ${GOPACKAGE}

type AuthorizationRepository interface {
	SetUserDB(ctx context.Context, user entity.User) (int, error)
	GetUserFromDB(ctx context.Context, user entity.User) (int, error)
	Validate(ctx context.Context, username string) error
}

type DataRepositoryInterface interface {
	SetRepData(id int, text []byte, eventType string) error
	GetRepData(id int, eventType string) ([]byte, error)
	DeleteRepData(id int, eventType string) error
}

type Repository struct {
	AuthorizationRepository
	DataRepositoryInterface
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthorizationRepository: NewAuthRepository(db),
		DataRepositoryInterface: NewDataRepository(db),
	}
}
