package dbx

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

//Client represent db client
type Client struct {
	*sqlx.DB
	transaction *Transaction
}

//ExecStatement create, update or update statement
func (me *Client) ExecStatement(statement *Statement) (sql.Result, error) {
	if me.transaction != nil {
		return me.transaction.ExecStatement(statement)
	}
	if me.DB == nil {
		return nil, errors.New("DB instance of type (*sql.DB) is nil")
	}
	return me.DB.NamedExecContext(statement.context, statement.SQL, statement.Parameters)
}

//Read records on database and return it as sqlx.Rows
func (me *Client) QueryStatement(statement *Statement) (*sqlx.Rows, error) {
	if me.transaction != nil && !me.transaction.IsComplete() {
		return me.transaction.QueryStatement(statement)
	}
	return me.DB.NamedQueryContext(statement.context, statement.SQL, statement.Parameters)
}

//BeginTransaction begin a new transaction
func (me *Client) BeginTransaction() (*Transaction, error) {
	sqlxTx, err := me.DB.Beginx()
	if err != nil {
		return nil, err
	}
	newTransaction := &Transaction{
		Tx: sqlxTx,
	}
	me.SetTransaction(newTransaction)
	return newTransaction, nil
}

//CompleteTransaction commit and reset current transaction
func (me *Client) CompleteTransaction() error {
	err := me.transaction.Commit()
	if err != nil {
		me.ResetTransaction()
		return err
	}

	me.ResetTransaction()

	return nil
}

//GetTransactionScope return current transaction
func (me *Client) GetTransaction() *Transaction {
	return me.transaction
}

//SetTransactionScope set the dbclient to which transaction it's belong to
func (me *Client) SetTransaction(transaction *Transaction) {
	me.transaction = transaction
}

//ResetTransaction set current transaction to nil
func (me *Client) ResetTransaction() {
	me.SetTransaction(nil)
}

//NewClient create new DB client instance
func NewClient(db *sqlx.DB) *Client {
	return &Client{
		DB: db,
	}
}