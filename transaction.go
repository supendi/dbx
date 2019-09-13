package dbx

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

//Transaction represent postgres transaction
type Transaction struct {
	*sqlx.Tx
}

//ExecuteCommand Create, Update or Delete statement
func (me *Transaction) ExecStatement(statement *Statement) (sql.Result, error) {
	return me.Tx.NamedExecContext(statement.context, statement.SQL, statement.Parameters)
}

//Read records on database and return it as sql.Rows
func (me *Transaction) QueryStatement(statement *Statement) (*sqlx.Rows, error) {
	return sqlx.NamedQueryContext(statement.context, me.Tx, statement.SQL, statement.Parameters)
}

//Commit the transaction
func (me *Transaction) Commit() error {
	return me.Tx.Commit()
}

//Rollback the transaction
func (me *Transaction) Rollback() error {
	return me.Tx.Rollback()
}

//NewTransaction create tx instance
func NewTransaction(db *sqlx.DB) (*Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	return &Transaction{
		Tx: tx,
	}, nil
}
