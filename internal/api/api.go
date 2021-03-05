package api

import (
	"github.com/PolyProjectOPD/Backend/internal/config"
	"github.com/PolyProjectOPD/Backend/internal/delivery/http"
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/internal/server"
	"github.com/PolyProjectOPD/Backend/internal/service"
	"github.com/PolyProjectOPD/Backend/pkg/hash"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(configPath string) {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.Init(configPath)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:  repos,
		Hasher: hasher,
	})
	handlers := http.NewHandler(services)

	srv := new(server.Server)
	if err = srv.Run(cfg, handlers.Init()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
