package migration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	//revive:disable:blank-imports
	_ "github.com/golang-migrate/migrate/v4/source/file"
	//revive:enable:blank-imports
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

const DefaultMigrationTable = "schema_migrations"

type Client struct {
	once     sync.Once
	migrator *migrate.Migrate
}

func New() *Client {
	return &Client{}
}

func (c *Client) Setup(ctx context.Context, db *sql.DB, path string, migrTable ...string) error {
	var (
		tableName = DefaultMigrationTable
		setupErr  error
	)

	if len(migrTable) > 0 && migrTable[0] != "" {
		tableName = migrTable[0]
	}

	c.once.Do(func() {
		post, err := c.postgres(ctx, db, tableName)
		if err != nil {
			setupErr = err

			return
		}

		filePath := fmt.Sprintf("file://%s", path)

		migr, err := migrate.NewWithDatabaseInstance(filePath, "postgres", post)
		if err != nil {
			setupErr = err

			return
		}

		c.migrator = migr
	})

	return setupErr
}

//revive:disable:cognitive-complexity
func (c *Client) SetupFS(ctx context.Context, db *sql.DB, filesystem fs.FS, migrTable ...string) error {
	var (
		tableName = DefaultMigrationTable
		setupErr  error
	)

	if len(migrTable) > 0 && migrTable[0] != "" {
		tableName = migrTable[0]
	}

	c.once.Do(func() {
		post, err := c.postgres(ctx, db, tableName)
		if err != nil {
			setupErr = err

			return
		}

		input, err := iofs.New(filesystem, ".")
		if err != nil {
			setupErr = err

			return
		}

		migr, err := migrate.NewWithInstance("iofs", input, "postgres", post)
		if err != nil {
			setupErr = err

			return
		}

		c.migrator = migr
	})

	return setupErr
}

//revive:enable:cognitive-complexity

func (c *Client) postgres(ctx context.Context, db *sql.DB, table string) (*postgres.Postgres, error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	cfg := &postgres.Config{
		MigrationsTable: table,
	}

	post, err := postgres.WithConnection(ctx, conn, cfg)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (c *Client) Up() error {
	if c.migrator == nil {
		return fmt.Errorf("migrate is not setup, use Setup() first")
	}

	err := c.migrator.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Down() error {
	if c.migrator == nil {
		return fmt.Errorf("migrate is not setup, use Setup() first")
	}

	return c.migrator.Down()
}
