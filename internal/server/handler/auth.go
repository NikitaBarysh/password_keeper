package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	_ "password_keeper/cmd/server/docs"
	"password_keeper/internal/common/entity"
	"password_keeper/internal/common/metrics"
	"password_keeper/internal/common/models"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce text
// @Param input body entity.User true "account info"
// @Success 201 {integer} integer 1
// @Failure 400 {string} err to parse body or err to unmarshal body
// @Failure 409 {string} not unique login or empty login
// @Failure 500 {string} err to create user or err to generate token
// @Router /register [post]

func (h *Handler) signUp(rw http.ResponseWriter, r *http.Request) {
	var statusCode int
	start := time.Now()
	defer func() {
		metrics.IncRequestStatus(r.URL.Path, statusCode)
		metrics.IncRequestDuration(r.URL.Path, time.Since(start))
	}()

	var input entity.User

	b, err := io.ReadAll(r.Body)
	if err != nil {
		statusCode = http.StatusBadRequest
		http.Error(rw, "err to parse body", statusCode)
		return
	}

	defer r.Body.Close()

	if err = json.Unmarshal(b, &input); err != nil {
		statusCode = http.StatusBadRequest
		http.Error(rw, "err to unmarshal body", statusCode)
		return
	}

	err = h.service.AuthorizationService.ValidateLogin(r.Context(), input)
	if errors.Is(err, models.ErrNotUniqueLogin) {
		statusCode = http.StatusConflict
		http.Error(rw, "not unique login or empty login", statusCode)
		return
	}

	id, err := h.service.AuthorizationService.CreateUser(r.Context(), input)
	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(rw, "err to create user", statusCode)
		return
	}

	token, err := h.service.AuthorizationService.GenerateJWTToken(id)
	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(rw, "err to generate token", statusCode)
		return
	}

	rw.Header().Set("Authorization", "Bearer "+token)
	statusCode = http.StatusCreated
	rw.WriteHeader(statusCode)
	rw.Write([]byte(token))
}

// @Summary SignIn
// @Tags auth
// @Description login account
// @ID login in account
// @Accept json
// @Produce text
// @Param input body entity.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400 {string} err to read body or err to unmarshal body
// @Failure 401 {string} invalid login or password
// @Failure 500 {string} err to generate token
// @Router /login [post]

func (h *Handler) singIn(rw http.ResponseWriter, r *http.Request) {
	var statusCode int
	start := time.Now()
	defer func() {
		metrics.IncRequestStatus(r.URL.Path, statusCode)
		metrics.IncRequestDuration(r.URL.Path, time.Since(start))
	}()

	var input entity.User

	b, err := io.ReadAll(r.Body)
	if err != nil {
		statusCode = http.StatusBadRequest
		http.Error(rw, "err to read body", statusCode)
		return
	}

	if err = json.Unmarshal(b, &input); err != nil {
		statusCode = http.StatusBadRequest
		http.Error(rw, "err to unmarshal body", statusCode)
		return
	}

	id, err := h.service.AuthorizationService.CheckData(r.Context(), input)
	if err != nil {
		statusCode = http.StatusUnauthorized
		http.Error(rw, "invalid login or password", statusCode)
		return
	}

	token, err := h.service.AuthorizationService.GenerateJWTToken(id)
	if err != nil {
		statusCode = http.StatusInternalServerError
		http.Error(rw, "err to generate token", statusCode)
		return
	}

	rw.Header().Set("Authorization", "Bearer "+token)
	statusCode = http.StatusOK
	rw.WriteHeader(statusCode)
	rw.Write([]byte(token))
}
