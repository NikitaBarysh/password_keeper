package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/common/models"
	"password_keeper/internal/server/service"
)

func TestSignUp(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string)

	tests := []struct {
		name               string
		inputBody          string
		inputUser          entity.User
		token              string
		mockBehaviour      mockBehaviour
		expectedStatusCode int
	}{
		{
			name:      "Test #1 sign-up successfully",
			inputBody: `{"login":"test", "password":"test"}`,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			token: "token",
			mockBehaviour: func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {
				s.EXPECT().ValidateLogin(ctx, user).Return(nil)
				s.EXPECT().CreateUser(ctx, user).Return(1, nil)
				s.EXPECT().GenerateJWTToken(1).Return(token, nil)
			},
			expectedStatusCode: 201,
		},
		{
			name:               "Test #2 sign-up bad input data",
			inputUser:          entity.User{},
			mockBehaviour:      func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Test #3 not unique login",
			inputBody: `{"login":"test", "password":"test"}`,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			token: "token",
			mockBehaviour: func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {
				s.EXPECT().ValidateLogin(ctx, user).Return(models.ErrNotUniqueLogin)
			},
			expectedStatusCode: 409,
		},
		{
			name:      "Test #5 err to create user",
			inputBody: `{"login":"test", "password":"test"}`,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			token: "token",
			mockBehaviour: func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {
				s.EXPECT().ValidateLogin(ctx, user).Return(nil)
				s.EXPECT().CreateUser(ctx, user).Return(0, errors.New("err to create user"))
			},
			expectedStatusCode: 500,
		},
		{
			name:      "Test #6 err to generate JWT",
			inputBody: `{"login":"test", "password":"test"}`,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			token: "token",
			mockBehaviour: func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {
				s.EXPECT().ValidateLogin(ctx, user).Return(nil)
				s.EXPECT().CreateUser(ctx, user).Return(1, nil)
				s.EXPECT().GenerateJWTToken(1).Return("", errors.New("err to generate JWT"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest(http.MethodPost, "/register",
				bytes.NewBufferString(test.inputBody))
			rw := httptest.NewRecorder()

			auth := service.NewMockAuthorizationService(c)
			test.mockBehaviour(req.Context(), auth, test.inputUser, test.token)

			serv := &service.Service{AuthorizationService: auth}

			handler := NewHandler(serv)

			handler.singUp(rw, req)

			res := rw.Result()
			assert.Equal(t, test.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestSignIn(t *testing.T) {
	type mockBehaviour func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string)
	tests := []struct {
		name               string
		inputBody          string
		inputUser          entity.User
		mockBehaviour      mockBehaviour
		token              string
		expectedStatusCode int
	}{
		{
			name:      "Test #1 sign-in successfully",
			inputBody: `{"login":"test", "password":"test"}`,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			mockBehaviour: func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {
				s.EXPECT().CheckData(ctx, user).Return(1, nil)
				s.EXPECT().GenerateJWTToken(1).Return(token, nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:      "Test #2 empty input data",
			inputBody: ``,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			mockBehaviour:      func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Test #3 no user with input data",
			inputBody: `{"login":"test", "password":"test"}`,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			mockBehaviour: func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {
				s.EXPECT().CheckData(ctx, user).Return(0, errors.New("err to check user"))
			},
			expectedStatusCode: 401,
		},
		{
			name:      "Test #4 err to generate JWT",
			inputBody: `{"login":"test", "password":"test"}`,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			mockBehaviour: func(ctx context.Context, s *service.MockAuthorizationService, user entity.User, token string) {
				s.EXPECT().CheckData(ctx, user).Return(1, nil)
				s.EXPECT().GenerateJWTToken(1).Return("", errors.New("err to generate JWT"))
			},
			expectedStatusCode: 500,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(test.inputBody))
			rw := httptest.NewRecorder()

			auth := service.NewMockAuthorizationService(c)

			test.mockBehaviour(req.Context(), auth, test.inputUser, test.token)

			serv := &service.Service{AuthorizationService: auth}

			handler := NewHandler(serv)

			handler.singIn(rw, req)

			res := rw.Result()
			assert.Equal(t, test.expectedStatusCode, res.StatusCode)
		})
	}
}
