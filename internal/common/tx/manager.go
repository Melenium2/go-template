package tx

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	once        sync.Once
	managerOnce *manager
)

type Func func(ctx context.Context) error

type Extension func(conn sqlx.ExtContext) sqlx.ExtContext

// ManagerTx is helper that provide possibility to execute group of database queries in
// single transaction.
type ManagerTx interface {
	// Do provide possibility to run transaction among all db functions inside txFunc.
	// Transaction level configurable by sql.Options. By default, will be applied
	// default level of used database, for example, 'read committed' in PostgreSQL.
	//
	// Example:
	//	func main() {
	//		tx.Manager().Do(ctx, func(ctx context.Context) error {
	//			res, err := repo.GetItem(ctx, 1)
	//			...
	//			err = nextRepo.SetItem(ctx, res)
	//		})
	//	}
	Do(ctx context.Context, txFunc Func, opts ...sql.TxOptions) error
	// Conn extracts transaction from provided context.Context. If transaction is not
	// contains inside context.Context, then function return default database connection.
	//
	// Example:
	//	func (r *repo) GetItem(ctx context.Context, id int64) (Item, error) {
	//		conn := tx.Manager().Conn(ctx)
	//
	//		res, err := conn.QueryContext(ctx, query)
	//		...
	//	}
	Conn(ctx context.Context) sqlx.ExtContext
}

type manager struct {
	db         *sqlx.DB
	decorators []Extension
}

func SetupManager(db *sqlx.DB, extensions ...Extension) {
	once.Do(func() {
		managerOnce = &manager{
			db:         db,
			decorators: extensions,
		}
	})
}

func Manager() ManagerTx {
	return managerOnce
}

func (m *manager) StartTx(ctx context.Context, opts ...sql.TxOptions) (context.Context, error) {
	if _, err := extractTx(ctx); err == nil {
		return ctx, nil
	}

	var txOpts sql.TxOptions

	if len(opts) > 0 {
		txOpts = opts[0]
	}

	tx, err := m.db.BeginTxx(ctx, &txOpts)
	if err != nil {
		return ctx, err
	}

	return injectTx(ctx, tx), nil
}

func (m *manager) Commit(ctx context.Context) error {
	tx, err := extractTx(ctx)
	if err != nil {
		return ErrTxNotFound
	}

	return tx.Commit()
}

func (m *manager) Rollback(ctx context.Context) error {
	tx, err := extractTx(ctx)
	if err != nil {
		return ErrTxNotFound
	}

	return tx.Rollback()
}

func (m *manager) Conn(ctx context.Context) sqlx.ExtContext {
	var conn sqlx.ExtContext = m.db

	tx, err := extractTx(ctx)
	if err == nil {
		conn = tx
	}

	return m.applyDecorators(conn)
}

func (m *manager) applyDecorators(conn sqlx.ExtContext) sqlx.ExtContext {
	for _, decorator := range m.decorators {
		conn = decorator(conn)
	}

	return conn
}

func (m *manager) Do(
	ctx context.Context,
	txFunc Func,
	opts ...sql.TxOptions,
) (err error) {
	_, err = extractTx(ctx)
	if err == nil {
		return txFunc(ctx)
	}

	txCtx, err := m.StartTx(ctx, opts...)
	if err != nil {
		return fmt.Errorf("can not start Tx, err, %w", err)
	}

	defer func() {
		// can easily panic here because of sqlx
		if p := recover(); p != nil {
			_ = m.Rollback(txCtx)

			debug.PrintStack()

			err = fmt.Errorf("recovered after panic in Tx.Do, err %v", p)
		}
	}()

	if err = txFunc(txCtx); err != nil {
		_ = m.Rollback(txCtx)

		return err
	}

	if err = m.Commit(txCtx); err != nil {
		return fmt.Errorf("error while committing Tx, err: %w", err)
	}

	return nil
}
