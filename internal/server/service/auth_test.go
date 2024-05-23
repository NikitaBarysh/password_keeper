package service

import (
	"context"
	"errors"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"password_keeper/config/server"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/server/repository"
)

const (
	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultDBUser     = "postgres"
	defaultDBPassword = "qwerty"
	defaultDBName     = "postgres"
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

const dbAddress = "postgres"

func TestServiceCreateUser(t *testing.T) {
	type mockBehaviour func(s *MockAuthorizationService)
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		user          entity.User
		wantErr       error
	}{
		{
			name:          "success",
			mockBehaviour: func(s *MockAuthorizationService) {},
			user: entity.User{
				Login:    "testSignUpService",
				Password: "test",
			},
			wantErr: nil,
		},
		{
			name:          "err to create",
			mockBehaviour: func(s *MockAuthorizationService) {},
			user: entity.User{
				Login:    "testSignUpService",
				Password: "test",
			},
			wantErr: errors.New("this login is busy"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := server.NewServConfig(
				server.WithDBAddress(getEnv("DB_HOST", defaultDBHost)),
				server.WithDBPort(getEnv("DB_PORT", defaultDBPort)),
				server.WithDBUsername(getEnv("DB_USER", defaultDBUser)),
				server.WithDBPassword(getEnv("DB_PASSWORD", defaultDBPassword)),
				server.WithDBDatabase(getEnv("DB_NAME", defaultDBName)),
			)

			db, err := repository.InitDataBase(ctx, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBUsername, cfg.DBPassword)
			require.NoError(t, err)

			defer db.Close()

			rep := repository.NewRepository(db)

			auth := NewMockAuthorizationService(c)
			test.mockBehaviour(auth)

			service := NewAuthService(rep, cfg)

			if _, err := service.CreateUser(ctx, test.user); (err != nil) != (test.wantErr != nil) {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, test.wantErr != nil)
			}

		})
	}
}

func TestServiceValidateLogin(t *testing.T) {
	type mockBehaviour func(s *MockAuthorizationService)
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		user          entity.User
		wantErr       error
	}{
		{
			name:          "success",
			mockBehaviour: func(s *MockAuthorizationService) {},
			user: entity.User{
				Login:    "someffjdnnnfmasdadsmssplks",
				Password: "test",
			},
			wantErr: nil,
		},
		{
			name:          "err to create",
			mockBehaviour: func(s *MockAuthorizationService) {},
			user: entity.User{
				Login:    "testSignUpService",
				Password: "test",
			},
			wantErr: errors.New("this login is busy"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := server.NewServConfig(
				server.WithDBAddress(dbAddress),
				server.WithDBPort("5432"),
				server.WithDBUsername("postgres"),
				server.WithDBPassword("qwerty"),
				server.WithDBDatabase("postgres"),
			)

			db, err := repository.InitDataBase(ctx, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBUsername, cfg.DBPassword)
			require.NoError(t, err)

			defer db.Close()

			rep := repository.NewRepository(db)

			auth := NewMockAuthorizationService(c)
			test.mockBehaviour(auth)

			service := NewAuthService(rep, cfg)

			if err := service.ValidateLogin(ctx, test.user); (err != nil) != (test.wantErr != nil) {
				t.Errorf("ValidateLogin() error = %v, wantErr %v", err, test.wantErr != nil)
			}

		})
	}
}

func TestServiceCheckData(t *testing.T) {
	type mockBehaviour func(s *MockAuthorizationService)
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		user          entity.User
		wantErr       error
	}{
		{
			name:          "success",
			mockBehaviour: func(s *MockAuthorizationService) {},
			user: entity.User{
				Login:    "testSignUpService",
				Password: "test",
			},
			wantErr: nil,
		},
		{
			name:          "wrong data",
			mockBehaviour: func(s *MockAuthorizationService) {},
			user: entity.User{
				Login:    "adminasdacasdc",
				Password: "adminacac",
			},
			wantErr: errors.New("wrong data for login"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := server.NewServConfig(
				server.WithDBAddress(dbAddress),
				server.WithDBPort("5432"),
				server.WithDBUsername("postgres"),
				server.WithDBPassword("qwerty"),
				server.WithDBDatabase("postgres"),
			)

			db, err := repository.InitDataBase(ctx, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBUsername, cfg.DBPassword)
			require.NoError(t, err)

			defer db.Close()

			rep := repository.NewRepository(db)

			auth := NewMockAuthorizationService(c)
			test.mockBehaviour(auth)

			service := NewAuthService(rep, cfg)

			if _, err := service.CheckData(ctx, test.user); (err != nil) != (test.wantErr != nil) {
				t.Errorf("CheckData() error = %v, wantErr %v", err, test.wantErr != nil)
			}

		})
	}
}

func TestGenerateJWT(t *testing.T) {
	type mockBehaviour func(s *MockAuthorizationService)
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		userID        int
		wantErr       error
	}{
		{
			name:          "success",
			mockBehaviour: func(s *MockAuthorizationService) {},
			userID:        1,
			wantErr:       nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg := server.NewServConfig(
				server.WithDBAddress(dbAddress),
				server.WithDBPort("5432"),
				server.WithDBUsername("postgres"),
				server.WithDBPassword("qwerty"),
				server.WithDBDatabase("postgres"),
			)

			db, err := repository.InitDataBase(ctx, cfg.DBHost, cfg.DBPort, cfg.DBDatabase, cfg.DBUsername, cfg.DBPassword)
			require.NoError(t, err)

			defer db.Close()

			rep := repository.NewRepository(db)

			auth := NewMockAuthorizationService(c)
			test.mockBehaviour(auth)

			service := NewAuthService(rep, cfg)

			if _, err := service.GenerateJWTToken(test.userID); (err != nil) != (test.wantErr != nil) {
				t.Errorf("GenerateJWT error = %v, wantErr %v", err, test.wantErr != nil)
			}
		})
	}
}
