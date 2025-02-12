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
		log.Fatalf("Error loading .env file")
	}
}

func main() {

	serverCfg := serverConfig.MustLoadServerConfig()
	storageCfg := storageConfig.MustLoadStorageConfig()
	accessToken.SetSecretKey(serverCfg.SecretKey)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil)) //TODO: check this out

	dataBase, err := storage.NewStorage(storageCfg)
	if err != nil {
		logger.Error("failed to connect to database", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())}) //TODO
		os.Exit(1)
	}

	logger.Info("connected to database")

	defer func() {
		err = dataBase.Close()
		if err != nil {
			logger.Error("got error when closing the DB connection", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			os.Exit(1)
		}
	}()

	repos := repository.NewRepository(dataBase)

	services := service.NewService(repos)

	handlers := handler.NewHandler(services)

	server := &http.Server{
		Addr:    serverCfg.Address,
		Handler: handlers.InitRoutes(),
	}

	go func() {
		slog.Info("starting server...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen and serve error", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

	log.Println("Server exiting")
}
