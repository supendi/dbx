package dbx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

//Transaction represent postgres transaction
type Transaction struct {
	*sqlx.Tx
	db         *sqlx.DB
	isComplete bool
}

//IsComplete determine if current transaction is already committed or rolledback
func (me *Transaction) IsComplete() bool {
	return me.isComplete
}

//ExecuteCommand Create, Update or Delete statement
func (me *Transaction) ExecStatementContext(ctx context.Context, statement *Statement) (sql.Result, error) {
	return me.Tx.NamedExecContext(ctx, statement.SQL, statement.Parameters)
}

//Read records on database and return it as sql.Rows
func (me *Transaction) QueryStatementContext(ctx context.Context, statement *Statement) (*sqlx.Rows, error) {
	return sqlx.NamedQueryContext(ctx, me.Tx, statement.SQL, statement.Parameters)
}

//ExecuteCommand Create, Update or Delete statement
func (me *Transaction) ExecStatement(statement *Statement) (sql.Result, error) {
	return me.Tx.NamedExec(statement.SQL, statement.Parameters)
}

//Read records on database and return it as sql.Rows
func (me *Transaction) QueryStatement(statement *Statement) (*sqlx.Rows, error) {
	return sqlx.NamedQuery(me.Tx, statement.SQL, statement.Parameters)
}

//Commit the transaction
func (me *Transaction) Commit() error {
	me.isComplete = true
	return me.Tx.Commit()
}

//Rollback the transaction
func (me *Transaction) Rollback() error {
	me.isComplete = true
	return me.Tx.Rollback()
}

//StartOver start over the transaction, current transaction will be overridden.
//it's recommended to always check if current transaction is complete or not
//by calling IsComplete() method
func (me *Transaction) StartOver() error {
	newTransaction, err := me.db.Beginx()
	if err != nil {
		return err
	}

	me.Tx = newTransaction
	me.isComplete = false
	return nil
}

//CommitAndStartOver commit current transaction and start a new one.
//If commit failed, it will try to rollback the transaction
func (me *Transaction) CommitAndStartOver() error {
	if me.Tx != nil && !me.IsComplete() {
		err := me.Commit()
		if err != nil {
			rollBackError := me.Rollback()
			if rollBackError != nil {
				return rollBackError
			}
			return err
		}
	}
	return me.StartOver()
}

//NewTransaction create tx instance
func NewTransaction(db *sqlx.DB) (*Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	return &Transaction{
		db: db,
		Tx: tx,
	}, nil
}
