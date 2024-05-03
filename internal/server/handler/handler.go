package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"password_keeper/internal/common/encryption"
	"password_keeper/internal/common/logger"
	"password_keeper/internal/server/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(newService *service.Service) *Handler {
	return &Handler{
		service: newService,
	}
}

func (h *Handler) Register(router *chi.Mux) {
	router.Route("/", func(router chi.Router) {
		router.Use(logger.LoggingMiddleware)

		router.Route("/", func(router chi.Router) {
			router.Use(encryption.DecryptMiddleware)
			router.Post("/register", h.singUp)
			router.Post("/login", h.singIn)
		})

		router.Route("/api", func(router chi.Router) {
			router.Use(h.AuthorizationMiddleware)
			router.Post("/set/{event}", h.setData)
			router.Get("/get/{event}", h.getData)
			router.Delete("/delete/{event}", h.deleteData)
		})
	})

	router.Route("/ws", func(router chi.Router) {
		router.Use(h.AuthorizationMiddleware)
		router.Handle("/connect", http.HandlerFunc(h.handleSetWebsocket))
	})

}
