package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	accessTokenTTL  = "auth.access_token_ttl"
	refreshTokenTTL = "auth.refresh_token_ttl"

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
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
	}
)

func authConfigInit() *AuthConfig {
	return &AuthConfig{
		JWTConfig: JWTConfig{
			SigningKey:      os.Getenv(jwtSigningKey),
			AccessTokenTTL:  viper.GetDuration(accessTokenTTL),
			RefreshTokenTTL: viper.GetDuration(refreshTokenTTL),
		},
		PasswordSalt: os.Getenv(passwordSalt),
	}
}
