package database

import (
	"fmt"

	"gihub.com/gibiw/api-example/internal/config"
	"github.com/jmoiron/sqlx"
)

func Initialize(cfg config.Database) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DatabaseName))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
