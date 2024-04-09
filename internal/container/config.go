package container

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Env string

var (
	Development Env = "dev"
	Stage       Env = "stage"
	Production  Env = "prod"
)

type Config struct {
	Environment Env    `env:"ENVIRONMENT" envDefault:"dev"`
	Branch      string `env:"BRANCH"`
	HTTPPort    string `env:"HTTP_PORT" envDefault:"4000"`

	DB   DB
	Amqp Amqp
}

type DB struct {
	Schema               string        `env:"DATABASE_SCHEMA" envDefault:"public"`
	Database             string        `env:"PGDATABASE" envDefault:"postgres"`
	Host                 string        `env:"PGHOST" envDefault:"localhost"`
	Port                 string        `env:"PGPORT" envDefault:"5432"`
	User                 string        `env:"PGUSER" evnDefault:"postgres"`
	Password             string        `env:"PGPASSWORD" envDefault:"postgres"`
	MaxOpenedConnections int           `env:"DATABASE_MAX_OPENED_CONNECTIONS" envDefault:"10"`
	MaxIdleTimeout       time.Duration `env:"DATABASE_MAX_IDLE_TIMEOUT" envDefault:"5m"`
}

type Amqp struct {
	AmqpHost     string `env:"AMQP_HOST" envDefault:"localhost"`
	AmqpVhost    string `env:"AMQP_VHOST" envDefault:"/"`
	AmqpPort     string `env:"AMQP_PORT" envDefault:"5672"`
	AmqpUser     string `env:"AMQP_USER" envDefault:"guest"`
	AmqpPassword string `env:"AMQP_PASSWORD" envDefault:"guest"`
}

func NewConfig() Config {
	// init envs from .env.example file
	_ = godotenv.Load()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	return cfg
}
