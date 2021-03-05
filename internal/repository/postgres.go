package repository

import (
	"fmt"
	"github.com/PolyProjectOPD/Backend/internal/config"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {

	dbConfig := cfg.DB
	db, err := sqlx.Connect(dbConfig.DriverName, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Name, dbConfig.Password, dbConfig.SSLMode))

	if err != nil {
		return nil, err
	}

	return db, nil
}
