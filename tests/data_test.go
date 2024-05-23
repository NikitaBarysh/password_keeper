package tests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"password_keeper/internal/common/entity"
)

func (s *APITestSuite) TestESetData() {
	r := s.Require()
	rctx := chi.NewRouter()
	s.handler.Register(rctx)

	login, password := "testingSetData", "testingSetData"

	pass := s.serv.GeneratePasswordHash(password)

	user := entity.User{Login: login, Password: pass}

	id, err := s.serv.CreateUser(context.Background(), user)
	if err != nil {
		s.Fail("Error checking data: " + err.Error())
	}

	jwt, err := s.serv.GenerateJWTToken(id)
	if err != nil {
		s.Fail("Error generating JWT: " + err.Error())
	}

	payload, testEvent := "test data", "testEvent"
	inputBody := fmt.Sprintf(`{"payload":"%s"}`, payload)

	req := httptest.NewRequest("POST", "/api/set/"+testEvent, bytes.NewReader([]byte(inputBody)))
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()

	rctx.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)
}

func (s *APITestSuite) TestFGetData() {
	r := s.Require()
	rctx := chi.NewRouter()
	s.handler.Register(rctx)

	login, password := "testingSetData", "testingSetData"

	pass := s.serv.GeneratePasswordHash(password)

	user := entity.User{Login: login, Password: pass}

	id, err := s.serv.CheckData(context.Background(), user)
	if err != nil {
		s.Fail("Error checking data: " + err.Error())
	}

	jwt, err := s.serv.GenerateJWTToken(id)
	if err != nil {
		s.Fail("Error generating JWT: " + err.Error())
	}

	testEvent := "testEvent"
	req := httptest.NewRequest("GET", "/api/get/"+testEvent, nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	res := httptest.NewRecorder()

	rctx.ServeHTTP(res, req)

	r.Equal(http.StatusOK, res.Result().StatusCode)
}

func (s *APITestSuite) TestGDeleteData() {
	r := s.Require()
	rctx := chi.NewRouter()
	s.handler.Register(rctx)

	login, password := "testingSetData", "testingSetData"
	pass := s.serv.GeneratePasswordHash(password)
	user := entity.User{Login: login, Password: pass}

	id, err := s.serv.CheckData(context.Background(), user)
	if err != nil {
		s.Fail("Error checking data: " + err.Error())
	}

	jwt, err := s.serv.GenerateJWTToken(id)
	if err != nil {
		s.Fail("Error generating JWT: " + err.Error())
	}

	testEvent := "testEvent"
	req := httptest.NewRequest("DELETE", "/api/delete/"+testEvent, nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	res := httptest.NewRecorder()

	rctx.ServeHTTP(res, req)

	r.Equal(http.StatusNoContent, res.Result().StatusCode)
}
