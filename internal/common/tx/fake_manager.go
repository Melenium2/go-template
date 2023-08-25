package tx

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func MockManager() {
	once.Do(func() {
		db, mock, _ := sqlmock.New()
		managerOnce = &manager{
			db: sqlx.NewDb(db, "postgres"),
		}

		mock.ExpectBegin()
		mock.ExpectCommit()
		mock.ExpectRollback()
	})
}
