package container

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/Melenium2/go-template/internal/common/tx"
	"github.com/Melenium2/go-template/pkg/migration"
	"github.com/Melenium2/go-template/pkg/psql"
)

func setupDatabase(cfg DB) *sqlx.DB {
	c := psql.Config{
		Schema:         cfg.Schema,
		DatabaseName:   cfg.Database,
		Host:           cfg.Host,
		Port:           cfg.Port,
		User:           cfg.User,
		Password:       cfg.Password,
		SimpleProtocol: true,
		Pool: psql.PoolConfig{
			MaxConnections:  cfg.MaxOpenedConnections,
			MinConnections:  1,
			MaxIddleTimeout: cfg.MaxIdleTimeout,
		},
	}

	conn, err := psql.Connect(c)
	if err != nil {
		log.Fatalf("error connecting database, %s", err.Error())
	}

	if err = conn.Ping(); err != nil {
		log.Fatalf("could not connect to database, %s", err.Error())
	}

	tx.SetupManager(conn)

	return conn
}

func setupMigrations(conn *sqlx.DB) {
	m := migration.New()

	err := m.Setup(context.TODO(), conn.DB, "db/migrations")
	if err != nil {
		log.Fatalf("can not setup migrations, err: %s", err)
	}

	if err = m.Up(); err != nil {
		log.Fatalf("error making migrations, %s", err)
	}
}

func makeDatabus(_ Amqp, _ Env, _ string) *Broker {
	return &Broker{}
}

func makeApps(_ *Container) *Apps {
	return &Apps{}
}

func makeClients(_ *Container, _ Config) *Clients {
	return &Clients{}
}

func makeStorages(_ *Container) *Storages {
	return &Storages{}
}

func makeServices(_ *Container) *Services {
	return &Services{}
}

func makeAppServices(_ *Container, _ Config) *ApplicationServices {
	return &ApplicationServices{}
}
