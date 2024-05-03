package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"password_keeper/internal/common/entity"
	"password_keeper/internal/common/models"
)

func (h *Handler) singUp(rw http.ResponseWriter, r *http.Request) {
	var input entity.User

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "err to parse body", http.StatusNotFound)
		return
	}

	defer r.Body.Close()

	if err = json.Unmarshal(b, &input); err != nil {
		http.Error(rw, "err to unmarshal body", http.StatusBadRequest)
		return
	}

	err = h.service.AuthorizationService.ValidateLogin(r.Context(), input)
	if errors.Is(err, models.ErrNotUniqueLogin) {
		http.Error(rw, "not unique login or empty login", http.StatusConflict)
		return
	}

	id, err := h.service.AuthorizationService.CreateUser(r.Context(), input)
	if err != nil {
		http.Error(rw, "err to create user", http.StatusInternalServerError)
		return
	}

	token, err := h.service.AuthorizationService.GenerateJWTToken(id)
	if err != nil {
		http.Error(rw, "err to generate token", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Authorization", "Bearer "+token)
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(token))
}

func (h *Handler) singIn(rw http.ResponseWriter, r *http.Request) {
	var input entity.User

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "err to read body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(b, &input); err != nil {
		http.Error(rw, "err to unmarshal body", http.StatusBadRequest)
		return
	}

	id, err := h.service.AuthorizationService.CheckData(r.Context(), input)
	if err != nil {
		http.Error(rw, "invalid login or password", http.StatusUnauthorized)
		return
	}

	token, err := h.service.AuthorizationService.GenerateJWTToken(id)
	if err != nil {
		http.Error(rw, "err to generate token", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Authorization", "Bearer "+token)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(token))
}
