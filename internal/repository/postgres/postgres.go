package postgres

import (
	"errors"
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/config"
	"github.com/jmoiron/sqlx"
)

const (
	UsersTable = "users"
)

func NewPostgresDB(cfg *config.DBConfig) (*sqlx.DB, error) {

	if cfg == nil {
		return nil, errors.New("empty config")
	}

	db, err := sqlx.Connect(cfg.DriverName, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	return db, nil
}
