package main

import (
	"Backend-trainee-assignment-winter-2025/internal/config/serverConfig"
	"Backend-trainee-assignment-winter-2025/internal/config/storageConfig"
	"Backend-trainee-assignment-winter-2025/internal/handler"
	accessToken "Backend-trainee-assignment-winter-2025/internal/jwt"
	"Backend-trainee-assignment-winter-2025/internal/repository"
	"Backend-trainee-assignment-winter-2025/internal/service"
	"Backend-trainee-assignment-winter-2025/internal/storage"
	"errors"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file")
		os.Exit(1)
	}
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	slog.SetDefault(logger)

	serverCfg, err := serverConfig.MustLoadServerConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	storageCfg, err := storageConfig.MustLoadStorageConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = accessToken.SetSecretKey(serverCfg.SecretKey)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = accessToken.SetTokenExpiry(serverCfg.TokenExpiry)
	if err != nil {
		log.Fatalf(err.Error())
	}

	dataBase, err := storage.NewStorage(storageCfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("connected to database")

	defer func() {
		err = dataBase.Close()
		if err != nil {
			slog.Error("failed to close database connection", "error", err)
			os.Exit(1)
		}
	}()

	repos := repository.NewRepository(dataBase)

	services := service.NewService(repos, logger)

	handlers := handler.NewHandler(services)

	server := &http.Server{
		Addr:    serverCfg.Address,
		Handler: handlers.InitRoutes(),
	}
	slog.Info(serverCfg.Address)

	go func() {
		slog.Info("starting server...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGABRT)
	<-quit
	slog.Info("shutting down server...")
	slog.Info("timeout of 5 seconds.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	select {
	case <-ctx.Done():
		slog.Info("timeout of 5 seconds reached")
	}

	slog.Info("server gracefully stopped")
}
