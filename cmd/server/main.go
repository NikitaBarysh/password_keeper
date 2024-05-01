package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"password_keeper/config/server"
	"password_keeper/internal/common/app"
	"password_keeper/internal/common/encryption"
	"password_keeper/internal/common/logger"
	"password_keeper/internal/server/handler"
	"password_keeper/internal/server/repository"
	"password_keeper/internal/server/service"
)

func main() {
	logging := logger.InitLogger()

	cfg := server.NewServer()
	logging.Info("Project cfg: endpoint: %s, database: %s", cfg.Endpoint, cfg.DataBaseDSN)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := repository.InitDataBase(ctx, cfg.DataBaseDSN, logging)
	if err != nil {
		logging.Fatal("get err from DB")
	}

	err = encryption.InitDecryptor(cfg.PrivateCryptoKeyPath)
	if err != nil {
		logging.Fatal("get err from encryption")
	}

	router := chi.NewRouter()

	newRepository := repository.NewRepository(db, logging)
	newService := service.NewService(newRepository, logging, cfg)
	newHandler := handler.NewHandler(newService, logging)

	newHandler.Register(router)

	srv := new(app.Server)
	go func() {
		if err := srv.Run(cfg.Endpoint, router); err != nil {
			logging.Fatal("Err to start server: ", err)
		}
	}()
	logging.Info("Server start")

	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGTERM, syscall.SIGINT)
	<-termSig

	if err = srv.ShutDown(ctx); err != nil {
		logging.Fatal("err to shutdown", err)
	}
}
