// Package repository - пакет в котором сохраняем данные в базу данных
package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"password_keeper/internal/common/entity"
)

//go:generate mockgen -source ${GOFILE} -destination mock.go -package ${GOPACKAGE}

// AuthorizationRepository - интерфейса с методами для работы с логином с паролем в базе
type AuthorizationRepository interface {
	SetUserDB(ctx context.Context, user entity.User) (int, error)
	GetUserFromDB(ctx context.Context, user entity.User) (int, error)
	Validate(ctx context.Context, username string) error
}

// DataRepositoryInterface - интерфейса с методами для работы с данными, который пользователь передал для хранения в базе
type DataRepositoryInterface interface {
	SetRepData(id int, text []byte, eventType string) error
	GetRepData(id int, eventType string) ([]byte, error)
	DeleteRepData(id int, eventType string) error
}

// Repository - хранит в себе интерфейсы пакет repository
type Repository struct {
	AuthorizationRepository
	DataRepositoryInterface
}

// NewRepository - создаем структуру Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthorizationRepository: NewAuthRepository(db),
		DataRepositoryInterface: NewDataRepository(db),
	}
}
