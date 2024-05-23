// Package app - пакет, где создаем сервер
package app

import (
	"context"
	"net/http"
)

// Server - структура, в которой находится сервер
type Server struct {
	httpServer *http.Server
}

// Run - запуск сервера с определенным адресом сервера и обработчиком
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    port,
		Handler: handler,
	}
	return s.httpServer.ListenAndServe()
}

// ShutDown - закрывает сервер
func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
