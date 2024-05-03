package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"password_keeper/config/server"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/common/models"
	"password_keeper/internal/server/repository"
)

const (
	tokenExp = time.Hour * 3
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

type AuthService struct {
	rep *repository.Repository
	cfg *server.ServConfig
}

func NewAuthService(newRep *repository.Repository, cfg *server.ServConfig) *AuthService {
	return &AuthService{
		rep: newRep,
		cfg: cfg,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user entity.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	id, err := s.rep.AuthorizationRepository.SetUserDB(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("err to get id from DB: %w", err)
	}
	return id, nil
}

func (s *AuthService) ValidateLogin(ctx context.Context, user entity.User) error {
	user.Password = s.generatePasswordHash(user.Password)
	id, err := s.rep.AuthorizationRepository.GetUserFromDB(ctx, user)
	if err != nil {
		return fmt.Errorf("err to get user: %w", err)
	}

	if id == 0 {
		return nil
	}

	return models.ErrNotUniqueLogin
}

func (s *AuthService) CheckData(ctx context.Context, user entity.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	id, err := s.rep.AuthorizationRepository.GetUserFromDB(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("err to get user fro DB: %w", err)
	}

	return id, nil
}

func (s *AuthService) GenerateJWTToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("err to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) GetUserIDFromToken(tokenString string) int {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}

	return claims.UserID
}

func (s *AuthService) generatePasswordHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return fmt.Sprintf("%x", h.Sum([]byte(s.cfg.Salt)))

}
