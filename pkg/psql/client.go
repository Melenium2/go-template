package psql

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	defaultMaxConnnections   = 10
	defaultMinConnections    = 1
	defaultMaxIddleTimeout   = 1 * time.Minute
	defaultHealthCheckPeriod = 30 * time.Second
)

// Конфигурация пула соединений клиента к постгресу.
type PoolConfig struct {
	// Кол-во максимальных соединений к серверу постгреса.
	//
	// Default: defaultMaxConnections.
	MaxConnections int
	// Минимальное кол-во активных соединений.
	//
	// Default: defaultMinConnections.
	MinConnections int
	// Максимальный таймаут простоя соединения. После
	// таймаута соединение убивается.
	//
	// Default: defaultMaxIddleTimeout.
	MaxIddleTimeout time.Duration
	// Таймаут на health check соединения.
	//
	// Default: defaultHealthCheckPeriod.
	HealthCheckPeriod time.Duration
}

// Config настройки подключения к базе данных.
type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	Schema       string
	// Насчет simple protocol можно почитать тут:
	// https://github.com/jackc/pgx/blob/master/conn.go#L627.
	SimpleProtocol bool
	// Конфигурация пула соединенеий.
	//
	// Optional.
	Pool PoolConfig
}

// URL создает connection URL в формате postgresql.
//
// Example:
//
//	postgres://jack:secret@pg.example.com:5432/mydb?pool_max_conns=10.
func (c Config) URL() string {
	urlstr := &url.URL{}
	urlstr.Scheme = "postgres"
	urlstr.User = url.UserPassword(c.User, c.Password)
	urlstr.Host = fmt.Sprintf("%s:%d", c.Host, c.Port)
	urlstr = urlstr.JoinPath(c.DatabaseName)

	pgProtocol := "cache_statement"
	if c.SimpleProtocol {
		pgProtocol = "simple_protocol"
	}

	query := urlstr.Query()
	query.Add("pool_max_conns", strconv.Itoa(c.Pool.MaxConnections))
	query.Add("pool_min_conns", strconv.Itoa(c.Pool.MinConnections))
	query.Add("pool_max_conn_idle_time", c.Pool.MaxIddleTimeout.String())
	query.Add("pool_health_check_period", c.Pool.HealthCheckPeriod.String())
	query.Add("default_query_exec_mode", pgProtocol)

	urlstr.RawQuery = query.Encode()

	return urlstr.String()
}

func defaultConfig() Config {
	port := os.Getenv("PGPORT")
	p, _ := strconv.Atoi(port)

	schema := "public"
	if customSchema := os.Getenv("PGSCHEMA"); customSchema != "" {
		schema = customSchema
	}

	return Config{
		Host:           os.Getenv("PGHOST"),
		Port:           p,
		User:           os.Getenv("PGUSER"),
		Password:       os.Getenv("PGPASSWORD"),
		DatabaseName:   os.Getenv("PGDATABASE"),
		Schema:         schema,
		SimpleProtocol: true,
		Pool: PoolConfig{
			MaxConnections:    defaultMaxConnnections,
			MinConnections:    defaultMinConnections,
			MaxIddleTimeout:   defaultMaxIddleTimeout,
			HealthCheckPeriod: defaultHealthCheckPeriod,
		},
	}
}

func mergeConfig(cfg1, cfg2 Config) Config {
	if cfg2.Host == "" {
		cfg2.Host = cfg1.Host
	}

	if cfg2.Port == 0 {
		cfg2.Port = cfg1.Port
	}

	if cfg2.User == "" {
		cfg2.User = cfg1.User
	}

	if cfg2.Password == "" {
		cfg2.Password = cfg1.Password
	}

	if cfg2.DatabaseName == "" {
		cfg2.DatabaseName = cfg1.DatabaseName
	}

	if cfg2.Schema == "" {
		cfg2.Schema = cfg1.Schema
	}

	if cfg2.Pool.MaxConnections == 0 {
		cfg2.Pool.MaxConnections = cfg1.Pool.MaxConnections
	}

	if cfg2.Pool.MinConnections == 0 {
		cfg2.Pool.MinConnections = cfg1.Pool.MinConnections
	}

	if cfg2.Pool.MaxIddleTimeout == 0 {
		cfg2.Pool.MaxIddleTimeout = cfg1.Pool.MaxIddleTimeout
	}

	if cfg2.Pool.HealthCheckPeriod == 0 {
		cfg2.Pool.HealthCheckPeriod = cfg1.Pool.HealthCheckPeriod
	}

	return cfg2
}

// Connect возвращает подключение к БД.
func Connect(c Config) (*sqlx.DB, error) {
	// Если нужна более тонкая конфигурация, то можно посмотреть тут, какие есть возможности.
	// https://github.com/jackc/pgx/blob/master/conn.go#L22.
	cfg := defaultConfig()
	cfg = mergeConfig(cfg, c)

	conf, err := pgxpool.ParseConfig(cfg.URL())
	if err != nil {
		return nil, err
	}

	conf.ConnConfig.RuntimeParams["search_path"] = c.Schema

	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	nativeConn := stdlib.OpenDBFromPool(pool)

	return sqlx.NewDb(nativeConn, "pgx"), nil
}
