package tests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"password_keeper/internal/common/encryption"
	"password_keeper/internal/common/entity"
)

func (s *APITestSuite) TestSignUp() {
	r := s.Require()

	rctx := chi.NewRouter()

	s.handler.Register(rctx)

	login, password := "testingSignUp", "testingSignUp"
	inputBody := fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password)

	enc, err := encryption.InitEncryptor(s.cfg.PublicCryptoKeyPath)
	if err != nil {
		s.FailNow("Failed to init encryptor", err)
	}

	b, err := enc.Encrypt([]byte(inputBody))
	if err != nil {
		s.FailNow("Failed to encrypt", err)
	}

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(b))

	resp := httptest.NewRecorder()

	rctx.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Code)

	pass := s.serv.AuthorizationService.GeneratePasswordHash(password)

	user := entity.User{
		Login:    login,
		Password: pass,
	}

	id, err := s.rep.GetUserFromDB(context.Background(), user)
	if err != nil {
		s.FailNow("Failed to get user", err)
	}

	r.Equal(2, id)

	err = s.rep.Validate(context.Background(), login)
	s.EqualError(err, "Validate: err to get id: %!w(<nil>) ")
}

func (s *APITestSuite) TestSignUpSameData() {
	r := s.Require()
	rctx := chi.NewRouter()
	s.handler.Register(rctx)

	login, password := "testSignIn", "testSignIn"
	inputBody := fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password)

	enc, err := encryption.InitEncryptor(s.cfg.PublicCryptoKeyPath)
	if err != nil {
		s.FailNow("Failed to init encryptor", err)
	}

	b, err := enc.Encrypt([]byte(inputBody))
	if err != nil {
		s.FailNow("Failed to encrypt", err)
	}

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(b))
	resp := httptest.NewRecorder()
	rctx.ServeHTTP(resp, req)

	r.Equal(http.StatusConflict, resp.Code)
}

func (s *APITestSuite) TestLogin() {
	r := s.Require()

	rctx := chi.NewRouter()

	s.handler.Register(rctx)

	login, password := "testSignIn", "testSignIn"
	inputBody := fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password)
	enc, err := encryption.InitEncryptor(s.cfg.PublicCryptoKeyPath)
	if err != nil {
		s.FailNow("Failed to init encryptor", err)
	}

	b, err := enc.Encrypt([]byte(inputBody))
	if err != nil {
		s.FailNow("Failed to encrypt", err)
	}

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(b))
	resp := httptest.NewRecorder()

	rctx.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Code)
}

func (s *APITestSuite) TestWrongLogin() {
	r := s.Require()

	rctx := chi.NewRouter()

	s.handler.Register(rctx)

	login, password := "wrongData", "wrongData"
	inputBody := fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password)

	enc, err := encryption.InitEncryptor(s.cfg.PublicCryptoKeyPath)
	if err != nil {
		s.FailNow("Failed to init encryptor", err)
	}

	b, err := enc.Encrypt([]byte(inputBody))
	if err != nil {
		s.FailNow("Failed to encrypt", err)
	}

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(b))
	resp := httptest.NewRecorder()

	rctx.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Code)
}
