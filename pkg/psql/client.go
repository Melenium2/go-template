package psql

import (
	"fmt"

	//revive:disable:blank-imports
	_ "github.com/jackc/pgx/v5/stdlib"
	//revive:enable:blank-imports
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Schema   string
	Database string
	Host     string
	Port     string
	User     string
	Password string
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database, c.Schema,
	)
}

func Connection(c Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", c.DSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}
