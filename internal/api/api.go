package api

import (
	"context"
	"github.com/PolyProjectOPD/Backend/internal/config"
	"github.com/PolyProjectOPD/Backend/internal/delivery/http"
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/internal/repository/postgres"
	"github.com/PolyProjectOPD/Backend/internal/server"
	"github.com/PolyProjectOPD/Backend/internal/service"
	"github.com/PolyProjectOPD/Backend/pkg/auth"
	"github.com/PolyProjectOPD/Backend/pkg/hash"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Network API
// @version 1.0
// @description API Server for Network OPD Project

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey UserAuth
// @in header
// @name Authorization
func Run(configPath string) {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.Init(configPath)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	db, err := postgres.NewPostgresDB(cfg.DB)
	if err != nil {
		logrus.Fatalf("failed to initialize postgres db: %s", err.Error())
	}

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)
	tokenManager, err := auth.NewManager(cfg.Auth.JWTConfig.SigningKey)
	if err != nil {
		logrus.Fatal("failed to initialize token manager: %s", err.Error())
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:           repos,
		Hasher:          hasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.JWTConfig.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.JWTConfig.RefreshTokenTTL,
	})
	handlers := http.NewHandler(services, tokenManager)

	srv := server.NewServer(cfg.HTTP, handlers.Init())
	go func() {
		if err = srv.Run(); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = time.Second * 5
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err = srv.Shutdown(ctx); err != nil {
		logrus.Errorf("failed to shutdown server: %v", err)
	}

	if err = db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
