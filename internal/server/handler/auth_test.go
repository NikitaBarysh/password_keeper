package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/common/logger"
	"password_keeper/internal/server/service"
)

func TestSignUp(t *testing.T) {
	type mockBehaviour func(s *service.MockAuthorizationService, user entity.User, token string)

	tests := []struct {
		name               string
		inputBody          string
		inputUser          entity.User
		token              string
		mockBehaviour      mockBehaviour
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "Test #1 sign-up successfully",
			inputBody: `{"username":"test", "password":"test"}`,
			inputUser: entity.User{
				Login:    "test",
				Password: "test",
			},
			token: "token",
			mockBehaviour: func(s *service.MockAuthorizationService, user entity.User, token string) {
				s.EXPECT().ValidateLogin(context.Background(), user.Login).Return(nil)
				s.EXPECT().CreateUser(context.Background(), user).Return(1, nil)
				s.EXPECT().GenerateJWTToken(1).Return(token, nil)
			},
			expectedStatusCode: 201,
			expectedBody:       `token`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logging := logger.InitLogger()
			c := gomock.NewController(t)
			defer c.Finish()

			auth := service.NewMockAuthorizationService(c)
			test.mockBehaviour(auth, test.inputUser, test.token)

			serv := &service.Service{AuthorizationService: auth}

			handler := NewHandler(serv, logging)
			req := httptest.NewRequest(http.MethodPost, "/register",
				bytes.NewBufferString(test.inputBody))
			rw := httptest.NewRecorder()
			handler.singUp(rw, req)

			res := rw.Result()
			assert.Equal(t, test.expectedStatusCode, res.StatusCode)
			assert.Equal(t, test.expectedBody, res.Body)
		})
	}
}
