package dbx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

//Transactioner TODO: must be added some more signatures
//Transactioner interface for dbclient
type Transactioner interface {
	ExecStatement(statement *Statement) (sql.Result, error)
	QueryStatement(statement *Statement) (*sqlx.Rows, error)
	ExecStatementContext(ctx context.Context, statement *Statement) (sql.Result, error)
	QueryStatementContext(ctx context.Context, statement *Statement) (*sqlx.Rows, error)
	Rollback() error
	Commit() error
}

//Context is dbclient too but wrap on or more statements.
type Context struct {
	Client
	Statements []*Statement
}

//AddStatement add new statement to context
func (me *Context) AddStatement(statement *Statement) {
	me.Statements = append(me.Statements, statement)
}

//AddStatements add statements to context
func (me *Context) AddStatements(statements ...*Statement) {
	for _, statement := range statements {
		me.Statements = append(me.Statements, statement)
	}
}

//ClearStatements clear current statements
func (me *Context) ClearStatements() {
	me.Statements = nil
}

//MustUseTransaction check if the context should use transaction or not
func (me *Context) MustUseTransaction() bool {
	return len(me.Statements) > 1 || me.transaction != nil
}

//execUseTransaction execute all deferred statements by using transaction
func (me *Context) execUseTransaction(ctx context.Context, transactioner Transactioner, statements []*Statement) ([]sql.Result, error) {
	var saveResults []sql.Result

	for _, statement := range statements {
		result, err := transactioner.ExecStatementContext(ctx, statement)
		if err != nil {
			if rollbackError := transactioner.Rollback(); rollbackError != nil {
				return nil, rollbackError
			}
			return nil, err
		}
		saveResults = append(saveResults, result)
	}
	if !me.IsUserDefinedTransaction {
		me.CompleteTransaction()
	}

	return saveResults, nil
}

//execWithoutTransaction execute statement without transaction
func (me *Context) execWithoutTransaction(ctx context.Context, statements []*Statement) ([]sql.Result, error) {
	var saveResults []sql.Result

	for _, statement := range statements {
		result, err := me.ExecStatementContext(ctx, statement)
		if err != nil {
			return nil, err
		}
		saveResults = append(saveResults, result)
	}

	return saveResults, nil
}

//SaveChanges execute all defered statements to database
func (me *Context) SaveChanges(ctx context.Context) ([]sql.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if me.MustUseTransaction() {
		if me.transaction == nil {

			me.IsUserDefinedTransaction = false //this flag is used in execUseTransaction(). if false, execUseTransaction will complete the transaction

			tx, err := me.Beginx()
			if err != nil {
				return nil, err
			}

			newTransaction := &Transaction{
				Tx: tx,
			}

			results, err := me.execUseTransaction(ctx, newTransaction, me.Statements)
			me.ClearStatements()
			if err != nil {
				return nil, err
			}

			err = newTransaction.Commit()
			if err != nil {
				return nil, err
			}
			return results, nil
		}
		results, err := me.execUseTransaction(ctx, me.transaction, me.Statements)
		me.ClearStatements()
		return results, err
	}

	results, err := me.execWithoutTransaction(ctx, me.Statements)
	me.ClearStatements()
	return results, err
}

//NewContext create new dbContext instance
func NewContext(dbClient *Client) *Context {
	return &Context{
		Client: *dbClient,
	}
}
