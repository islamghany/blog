package db

import (
	"fmt"
	"github/islamghany/blog/business/data/transaction"

	"github.com/jmoiron/sqlx"
)

type DBBeginner struct {
	sqlxDB *sqlx.DB
}

// NewBeginner constructs a value that implements the beginner interface.
func NewBeginner(sqlxDB *sqlx.DB) *DBBeginner {
	return &DBBeginner{
		sqlxDB: sqlxDB,
	}
}

func (db *DBBeginner) Begin() (transaction.Transaction, error) {
	return db.sqlxDB.Beginx()
}

// GetExtContext is a helper function that extracts the sqlx value
// from the domain transactor interface for transactional use.
func GetExtContext(tx transaction.Transaction) (sqlx.ExtContext, error) {
	ec, ok := tx.(sqlx.ExtContext)
	if !ok {
		return nil, fmt.Errorf("Transactor(%T) not of a type *sql.Tx", tx)
	}

	return ec, nil
}
