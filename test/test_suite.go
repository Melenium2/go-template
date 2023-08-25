package test

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/Melenium2/go-tempalte/internal/common/tx"
)

type config struct {
	Schema   string `env:"DATABASE_SCHEMA" envDefault:"public"`
	Database string `env:"PGDATABASE" envDefault:"postgres"`
	Host     string `env:"PGHOST" envDefault:"localhost"`
	Port     string `env:"PGPORT" envDefault:"5432"`
	User     string `env:"PGUSER" evnDefault:"postgres"`
	Password string `env:"PGPASSWORD" envDefault:"postgres"`
}

func newConfig() config {
	envPath := path2env()
	// init envs from .env file
	_ = godotenv.Load(envPath)

	var cfg config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	return cfg
}

func path2env() string {
	_, c, _, _ := runtime.Caller(0)
	s := path.Join(path.Dir(c))
	sp := strings.Split(s, string(os.PathSeparator))
	sp = sp[:len(sp)-1]
	sp = append([]string{string(os.PathSeparator)}, sp...)

	return path.Join(append(sp, ".env")...)
}

type DBSuite struct {
	suite.Suite

	Conn *sqlx.DB
}

func NewSuite() DBSuite {
	return DBSuite{}
}

func (suite *DBSuite) SetupSuite() error {
	cfg := newConfig()

	var (
		address = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
		)
	)

	db, err := sqlx.Connect("postgres", address)
	if err != nil {
		return fmt.Errorf("postgres is not connected, run local instance of postgres, err: %w", err)
	}

	suite.Conn = db

	tx.SetupManager(suite.Conn)

	return nil
}
