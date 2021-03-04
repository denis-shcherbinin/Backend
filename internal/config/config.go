package config

import "time"

const (
	defaultHttpPort               = "8080"
	defaultHttpRWTimeout          = time.Second * 10
	defaultHttpMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = time.Minute * 15
	defaultRefreshTokenTTL        = time.Hour * 24 * 30 // 30 days
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
		Host               string
		Port               string
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxHeaderMegabytes int
	}

	AuthConfig struct {
		JWT  JWTConfig
		salt string
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
	}
)
