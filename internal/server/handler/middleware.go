package handler

import (
	"context"
	"net/http"
	"strings"
)

const authorizationHeader = "Authorization"

func (h *Handler) AuthorizationMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			http.Error(rw, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(rw, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if headerParts[1] == "" {
			http.Error(rw, "empty token", http.StatusUnauthorized)
			return
		}

		userID := h.service.AuthorizationService.GetUserIDFromToken(headerParts[1])
		if userID == -1 {
			http.Error(rw, "err to parse token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", userID)
		handler.ServeHTTP(rw, r.WithContext(ctx))
	})
}
