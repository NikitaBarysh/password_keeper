package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"password_keeper/internal/common/metrics"

	"github.com/go-chi/chi/v5"
	"password_keeper/config/server"
	"password_keeper/internal/common/app"
	"password_keeper/internal/common/encryption"
	"password_keeper/internal/common/logger"
	"password_keeper/internal/server/handler"
	"password_keeper/internal/server/repository"
	"password_keeper/internal/server/service"
)

// @Title Password keeper app API
// @Version 1.0
// @Description App, which save secret data

// @Host localhost:8000

// @SecurityDefinitions.apiKey ApiKeyAuth
// @in header
// @Name Authorization

func main() {
	logging := logger.InitLogger()

	cfg := server.NewServer()
	logging.Info(fmt.Sprintf("Project cfg: endpoint: %s", cfg.Endpoint))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := repository.InitDataBase(ctx,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
		cfg.DBUsername,
		cfg.DBPassword,
	)
	if err != nil {
		logging.Fatal("Main client: %s", err.Error())
	}

	err = encryption.InitDecryptor(cfg.PrivateCryptoKeyPath)
	if err != nil {
		logging.Fatal("Main client: %s", err.Error())
	}

	router := chi.NewRouter()

	newRepository := repository.NewRepository(db)
	newService := service.NewService(newRepository, cfg)
	newHandler := handler.NewHandler(newService)

	newHandler.Register(router)

	go func() {
		_ = metrics.Listen("localhost:8082")
	}()
	logging.Info("Metric server listening on localhost:8082/metrics")

	srv := new(app.Server)
	go func() {
		if err = srv.Run(cfg.Endpoint, router); err != nil {
			logging.Fatal("Main server:Err to start client: %s", err.Error())
		}
	}()
	logging.Info("Server start")

	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGTERM, syscall.SIGINT)
	<-termSig

	if err = srv.ShutDown(ctx); err != nil {
		logging.Fatal("err to shutdown %s", err.Error())
	}
}
