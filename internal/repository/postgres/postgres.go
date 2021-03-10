package postgres

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/config"
	"github.com/jmoiron/sqlx"
)

const (
	UsersTable = "users"
	UsersSessionsTable = "users_sessions"
)

func NewPostgresDB(cfg *config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect(cfg.DriverName, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	return db, nil
}
