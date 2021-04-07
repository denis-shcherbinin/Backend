package config

import (
	"github.com/spf13/viper"
	"os"
)

const (
	dbName       = "db.name"
	dbHost       = "db.host"
	dbPort       = "db.port"
	dbUser       = "db.user"
	dbSSLMode    = "db.sslmode"
	dbDriverName = "db.driver_name"

	dbPassword = "POSTGRES_PASSWORD"
)

type DBConfig struct {
	Name       string
	Host       string
	Port       string
	User       string
	SSLMode    string
	Password   string
	DriverName string
}

func dbConfigInit() *DBConfig {
	return &DBConfig{
		Name:       viper.GetString(dbName),
		Host:       viper.GetString(dbHost),
		Port:       viper.GetString(dbPort),
		User:       viper.GetString(dbUser),
		SSLMode:    viper.GetString(dbSSLMode),
		Password:   os.Getenv(dbPassword),
		DriverName: viper.GetString(dbDriverName),
	}
}
