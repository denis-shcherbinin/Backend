package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	DB   *DBConfig
	HTTP *HTTPConfig
	Auth *AuthConfig
	Storage *StorageConfig
}

func Init(configPath string) (*Config, error) {

	path := strings.Split(configPath, "/")
	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New(fmt.Sprintf("error initializing configs: %s", err.Error()))
	}

	if err := godotenv.Load(); err != nil {
		return nil, errors.New(fmt.Sprintf("error loading env variables: %s", err.Error()))
	}

	return &Config{
		DB:   dbConfigInit(),
		HTTP: httpConfigInit(),
		Auth: authConfigInit(),
		Storage: storageConfigInit(),
	}, nil
}
