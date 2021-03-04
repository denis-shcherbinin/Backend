package app

import (
	"github.com/PolyProjectOPD/Backend/internal/delivery/http"
	"github.com/PolyProjectOPD/Backend/internal/repository"
	"github.com/PolyProjectOPD/Backend/internal/server"
	"github.com/PolyProjectOPD/Backend/internal/service"
	"github.com/PolyProjectOPD/Backend/pkg/hash"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func Run() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		DriverName: viper.GetString("db.driver_name"),
		Host:       viper.GetString("db.host"),
		Port:       viper.GetString("db.port"),
		User:       viper.GetString("db.user"),
		DBName:     viper.GetString("db.dbname"),
		SSLMode:    viper.GetString("db.sslmode"),
		Password:   os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	hasher := hash.NewSHA1Hasher(os.Getenv("PASSWORD_SALT"))

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:  repos,
		Hasher: hasher,
	})
	handlers := http.NewHandler(services)

	srv := new(server.Server)
	if err = srv.Run(viper.GetString("http.port"), handlers.Init()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
