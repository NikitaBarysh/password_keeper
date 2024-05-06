// Package service - пакет в котором реализовано безнес-логика сервера
package service

import (
	"context"

	"password_keeper/config/server"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/server/repository"
)

//go:generate mockgen -source ${GOFILE} -destination mock.go -package ${GOPACKAGE}

// Service - структура в которой хранится интерфейсы пакета service
type Service struct {
	AuthorizationService
	DataServiceInterface
}

// AuthorizationService - интерфейс с методами аутентификации и регистрации
type AuthorizationService interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	ValidateLogin(ctx context.Context, user entity.User) error
	CheckData(ctx context.Context, user entity.User) (int, error)
	GenerateJWTToken(userID int) (string, error)
	GetUserIDFromToken(tokenString string) int
}

// DataServiceInterface - интерфейс с методами добавления, удаления, получения данных из базы
type DataServiceInterface interface {
	SetData(id int, text []byte, eventType string) error
	GetData(id int, eventType string) ([]byte, error)
	DeleteData(id int, eventType string) error
}

// NewService - создаем структуру Service
func NewService(rep *repository.Repository, cfg *server.ServConfig) *Service {
	return &Service{
		AuthorizationService: NewAuthService(rep, cfg),
		DataServiceInterface: NewDataService(rep),
	}
}
