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
	logging *logger.Logger
}

func NewHandler(newService *service.Service, log *logger.Logger) *Handler {
	return &Handler{
		service: newService,
		logging: log,
	}
}

func (h *Handler) Register(router *chi.Mux) {
	router.Use(logger.LoggingMiddleware) // Не использую потому что не могу тогда создать websocket соединение
	// из-за http.Hijacker, попробовал способ (см. файл с midl.) не получается

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
		router.Handle("/ws", http.HandlerFunc(h.handleSetWebsocket))
	})
}
