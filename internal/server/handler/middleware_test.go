package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"password_keeper/internal/server/service"
)

func TestAuthMiddleware(t *testing.T) {
	type mockBehaviour func(s *service.MockAuthorizationService, token string)
	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		expectedStatus int
		token          string
		headerName     string
		headerValue    string
	}{
		{
			name: "OK",
			mockBehaviour: func(s *service.MockAuthorizationService, token string) {
				s.EXPECT().GetUserIDFromToken(token).Return(1)
			},
			token:          "test",
			expectedStatus: http.StatusOK,
			headerName:     "Authorization",
			headerValue:    "Bearer test",
		},
		{
			name:           "No header",
			mockBehaviour:  func(s *service.MockAuthorizationService, token string) {},
			token:          "test",
			headerName:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			mockBehaviour:  func(s *service.MockAuthorizationService, token string) {},
			token:          "test",
			headerName:     "Authorization",
			headerValue:    "Bearer",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token value",
			mockBehaviour:  func(s *service.MockAuthorizationService, token string) {},
			token:          "test",
			headerName:     "Authorization",
			headerValue:    "Bearer ",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "err to parse token",
			mockBehaviour: func(s *service.MockAuthorizationService, token string) {
				s.EXPECT().GetUserIDFromToken(token).Return(-1)
			},
			token:          "test",
			headerName:     "Authorization",
			headerValue:    "Bearer test",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := service.NewMockAuthorizationService(c)
			test.mockBehaviour(auth, test.token)

			serv := &service.Service{AuthorizationService: auth}
			handler := NewHandler(serv)

			rcxt := chi.NewRouter()
			rcxt.Use(handler.AuthorizationMiddleware)

			rcxt.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
				id, _ := r.Context().Value("user").(int)
				str := strconv.Itoa(id)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(str))
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			w := httptest.NewRecorder()
			req.Header.Set(test.headerName, test.headerValue)

			rcxt.ServeHTTP(w, req)
		})
	}
}
