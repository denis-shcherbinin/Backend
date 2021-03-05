package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type (
	Config struct {
		DB   DBConfig
		HTTP HTTPConfig
		Auth AuthConfig
	}

	DBConfig struct {
		Name       string
		Host       string
		Port       string
		User       string
		SSLMode    string
		Password   string
		DriverName string
	}

	HTTPConfig struct {
		Port               string
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxHeaderMegabytes int
	}

	AuthConfig struct {
		JWTConfig    JWTConfig
		PasswordSalt string
	}

	JWTConfig struct {
		SigningKey      string
		accessTokenTTL  time.Duration
		refreshTokenTTL time.Duration
	}
)

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
		DB: DBConfig{
			Name:       viper.GetString("db.name"),
			Host:       viper.GetString("db.host"),
			Port:       viper.GetString("db.port"),
			User:       viper.GetString("db.user"),
			SSLMode:    viper.GetString("db.sslmode"),
			Password:   os.Getenv("DB_PASSWORD"),
			DriverName: viper.GetString("db.driver_name"),
		},
		HTTP: HTTPConfig{
			Port:               viper.GetString("http.port"),
			ReadTimeout:        viper.GetDuration("read_timeout"),
			WriteTimeout:       viper.GetDuration("write_timeout"),
			MaxHeaderMegabytes: viper.GetInt("http.max_header_bytes"),
		},
		Auth: AuthConfig{
			JWTConfig: JWTConfig{
				SigningKey:      os.Getenv("JWT_SIGNING_KEY"),
				accessTokenTTL:  viper.GetDuration("access_token_ttl"),
				refreshTokenTTL: viper.GetDuration("refresh_token_ttl"),
			},
			PasswordSalt: os.Getenv("PASSWORD_SALT"),
		},
	}, nil
}
