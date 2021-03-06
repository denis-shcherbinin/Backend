package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	accessTokenTTL  = "access_token_ttl"
	refreshTokenTTL = "refresh_token_ttl"

	jwtSigningKey = "JWT_SIGNING_KEY"
	passwordSalt  = "PASSWORD_SALT"
)

type (
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

func authConfigInit() *AuthConfig {
	return &AuthConfig{
		JWTConfig: JWTConfig{
			SigningKey:      os.Getenv(jwtSigningKey),
			accessTokenTTL:  viper.GetDuration(accessTokenTTL),
			refreshTokenTTL: viper.GetDuration(refreshTokenTTL),
		},
		PasswordSalt: os.Getenv(passwordSalt),
	}
}
