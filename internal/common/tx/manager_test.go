package tx

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type ManagerSuite struct {
	suite.Suite

	db      *sqlx.DB
	sqlMock sqlmock.Sqlmock
}

func (suite *ManagerSuite) SetupSuite() {
	db, mock, _ := sqlmock.New()

	suite.sqlMock = mock
	suite.db = sqlx.NewDb(db, "postgres")
	SetupManager(suite.db)
}

func (suite *ManagerSuite) TestStartTx_Should_start_new_transaction_and_set_transaction_pointer_to_context() {
	ctx := context.Background()

	suite.sqlMock.ExpectBegin()

	tx, err := managerOnce.StartTx(ctx)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(tx)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *ManagerSuite) TestStartTx_Should_return_same_context_as_passed_because_tx_already_begin() {
	expected := &sqlx.Tx{}
	ctx := injectTx(context.Background(), expected)

	ctx, err := managerOnce.StartTx(ctx)
	suite.Assert().NoError(err)

	tx, err := extractTx(ctx)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, tx)
}

func (suite *ManagerSuite) TestStartTx_Should_return_error_if_can_not_begin_tx() {
	suite.sqlMock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	_, err := managerOnce.StartTx(context.Background())
	suite.Assert().ErrorIs(err, sql.ErrConnDone)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *ManagerSuite) TestCommit_Should_commit_changes_without_error() {
	suite.sqlMock.ExpectBegin()

	ctx, _ := managerOnce.StartTx(context.Background())

	suite.sqlMock.ExpectCommit()

	err := managerOnce.Commit(ctx)
	suite.Assert().NoError(err)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *ManagerSuite) TestCommit_Should_return_error_if_context_does_not_have_tx() {
	err := managerOnce.Commit(context.Background())
	suite.Assert().ErrorIs(err, ErrTxNotFound)
}

func (suite *ManagerSuite) TestRollback_Should_rollback_changes_without_error() {
	suite.sqlMock.ExpectBegin()

	ctx, _ := managerOnce.StartTx(context.Background())

	suite.sqlMock.ExpectRollback()

	err := managerOnce.Rollback(ctx)
	suite.Assert().NoError(err)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *ManagerSuite) TestRollback_Should_return_error_if_context_does_not_have_tx() {
	err := managerOnce.Rollback(context.Background())
	suite.Assert().ErrorIs(err, ErrTxNotFound)
}

func (suite *ManagerSuite) TestConn_Should_return_tx_from_context() {
	tx := &sqlx.Tx{}
	ctx := injectTx(context.Background(), tx)

	conn := Manager().Conn(ctx)
	suite.Assert().Equal(tx, conn)
}

func (suite *ManagerSuite) TestConn_Should_return_default_conn() {
	conn := Manager().Conn(context.Background())
	suite.Assert().Equal(suite.db, conn)
}

func (suite *ManagerSuite) TestDo_Should_run_txFunc_with_new_transaction() {
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	err := Manager().Do(context.Background(), func(ctx context.Context) error {
		a := 1 + 1
		_ = a

		return nil
	})
	suite.Assert().NoError(err)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *ManagerSuite) TestDo_Should_return_error_in_txFunc_and_rollback_tx() {
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	err := Manager().Do(context.Background(), func(ctx context.Context) error {
		return errors.New("error")
	})
	suite.Assert().Error(err)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *ManagerSuite) TestDo_Should_return_error_if_can_not_commit_changes() {
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit().WillReturnError(sql.ErrTxDone)

	err := Manager().Do(context.Background(), func(ctx context.Context) error {
		a := 2 * 2
		_ = a

		return nil
	})
	suite.Assert().ErrorIs(err, sql.ErrTxDone)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *ManagerSuite) TestDo_Should_panic_in_txFunc_then_recover_and_rollback() {
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	err := Manager().Do(context.Background(), func(ctx context.Context) error {
		panic("panic a!a!a!")

		return nil
	})
	suite.Assert().Error(err)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func (suite *ManagerSuite) TestDo_Should_rollback_if_context_was_canceled() {
	suite.T().Parallel()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit().WillReturnError(sql.ErrTxDone)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := Manager().Do(ctx, func(ctx context.Context) error {
		time.Sleep(150 * time.Millisecond)

		return nil
	})
	suite.Assert().ErrorIs(err, sql.ErrTxDone)
}

func (suite *ManagerSuite) TestDo_Should_run_single_transaction_among_some_Do_functions() {
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	err := Manager().Do(context.Background(), func(ctx context.Context) error {
		err := Manager().Do(ctx, func(ctx context.Context) error {
			err := Manager().Do(ctx, func(ctx context.Context) error {
				return nil
			})

			return err
		})

		return err
	})
	suite.Assert().NoError(err)
	suite.Assert().NoError(suite.sqlMock.ExpectationsWereMet())
}

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}
