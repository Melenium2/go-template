package test

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/Melenium2/go-template/internal/common/tx"
	"github.com/Melenium2/go-template/pkg/psql"
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

	port, _ := strconv.Atoi(cfg.Port)

	psqlConfig := psql.Config{
		Host:           cfg.Host,
		Port:           port,
		User:           cfg.User,
		Password:       cfg.Password,
		DatabaseName:   cfg.Database,
		Schema:         cfg.Schema,
		SimpleProtocol: true,
	}

	db, err := psql.Connect(psqlConfig)
	if err != nil {
		return fmt.Errorf("postgres connection is not established, %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("postgres is not connected, run local instance of postgres, %w", err)
	}

	suite.Conn = db

	tx.SetupManager(suite.Conn)

	return nil
}
