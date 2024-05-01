package service

import (
	"context"

	"password_keeper/config/server"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/common/logger"
	"password_keeper/internal/server/repository"
)

//go:generate mockgen -source ${GOFILE} -destination mock.go -package ${GOPACKAGE}

type Service struct {
	AuthorizationService
	DataServiceInterface
}

type AuthorizationService interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	ValidateLogin(ctx context.Context, user entity.User) error
	CheckData(ctx context.Context, user entity.User) (int, error)
	GenerateJWTToken(userID int) (string, error)
	GetUserIDFromToken(tokenString string) int
}

type DataServiceInterface interface {
	SetData(ctx context.Context, id int, text []byte, eventType string) error
	GetData(ctx context.Context, id int, eventType string) ([]byte, error)
	DeleteData(ctx context.Context, id int, eventType string) error
}

func NewService(rep *repository.Repository, log *logger.Logger, cfg *server.ServConfig) *Service {
	return &Service{
		AuthorizationService: NewAuthService(rep, log, cfg),
		DataServiceInterface: NewDataService(rep, log),
	}
}
