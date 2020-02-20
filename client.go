package dbx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

//Client represent db client
type Client struct {
	*sqlx.DB
	transaction              *Transaction
	IsUserDefinedTransaction bool
}

//ExecStatement create, update or update statement
func (me *Client) ExecStatement(statement *Statement) (sql.Result, error) {
	if me.transaction != nil {
		return me.transaction.ExecStatement(statement)
	}
	if me.DB == nil {
		return nil, errors.New("DB instance of type (*sql.DB) is nil")
	}
	return me.DB.NamedExec(statement.SQL, statement.Parameters)
}

//ExecStatementContext create, update or update statement
func (me *Client) ExecStatementContext(ctx context.Context, statement *Statement) (sql.Result, error) {
	if me.transaction != nil {
		return me.transaction.ExecStatementContext(ctx, statement)
	}
	if me.DB == nil {
		return nil, errors.New("DB instance of type (*sql.DB) is nil")
	}
	return me.DB.NamedExecContext(ctx, statement.SQL, statement.Parameters)
}

//QueryStatement records on database and return it as sqlx.Rows
func (me *Client) QueryStatement(statement *Statement) (*sqlx.Rows, error) {
	if me.transaction != nil && !me.transaction.IsComplete() {
		return me.transaction.QueryStatement(statement)
	}
	return me.DB.NamedQuery(statement.SQL, statement.Parameters)
}

//QueryStatementContext records on database and return it as sqlx.Rows
func (me *Client) QueryStatementContext(ctx context.Context, statement *Statement) (*sqlx.Rows, error) {
	if me.transaction != nil && !me.transaction.IsComplete() {
		return me.transaction.QueryStatementContext(ctx, statement)
	}
	return me.DB.NamedQueryContext(ctx, statement.SQL, statement.Parameters)
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
	if me.transaction != nil && !me.transaction.isComplete {
		err := me.transaction.Commit()
		if err != nil {
			if rollbackError := me.transaction.Rollback(); rollbackError != nil {
				return rollbackError
			}
			me.ResetTransaction()
			return err
		}
	}

	me.ResetTransaction()

	return nil
}

//GetTransaction return current transaction
func (me *Client) GetTransaction() *Transaction {
	return me.transaction
}

//SetTransaction set the dbclient to which transaction it's belong to
func (me *Client) SetTransaction(transaction *Transaction) {
	me.transaction = transaction
	if transaction != nil {
		me.IsUserDefinedTransaction = true
	}
}

//SetTransactionScope set transaction to the specified context
func (me *Context) SetTransactionScope(ctx context.Context) error {
	newTransaction, err := NewTransaction(me.DB)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, contextKey, newTransaction)
	return nil
}

//GetTransactionScope get transaction from context
func (me *Context) GetTransactionScope(ctx context.Context) *Transaction {
	value := ctx.Value(contextKey)
	if value != nil {
		return value.(*Transaction)
	}
	return nil
}

//ResetTransaction set current transaction to nil
func (me *Client) ResetTransaction() {
	me.IsUserDefinedTransaction = false
	me.SetTransaction(nil)
}

//NewClient create new DB client instance
func NewClient(db *sqlx.DB) *Client {
	return &Client{
		DB: db,
	}
}
