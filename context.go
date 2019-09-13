package dbx

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

//TODO: must be added some more signatures
//Execer interface for dbclient
type Execer interface {
	ExecStatement(statement *Statement) (sql.Result, error)
	QueryStatement(statement *Statement) (*sqlx.Rows, error)
}

//TODO: must be added some more signatures
//Execer interface for dbclient
type Transactioner interface {
	ExecStatement(statement *Statement) (sql.Result, error)
	QueryStatement(statement *Statement) (*sqlx.Rows, error)
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

//AddStatement add statements to context
func (me *Context) AddStatements(statements ...*Statement) {
	for _, statement := range statements {
		me.Statements = append(me.Statements, statement)
	}
}

//ClearStatements clear current statements
func (me *Context) ClearStatements() {
	me.Statements = nil
}

//ShouldUseTransaction check if the context should use transaction or not
func (me *Context) ShouldUseTransaction() bool {
	return len(me.Statements) > 1 || me.transaction != nil
}

//execUseTransaction execute all deferred statements by using transaction
func (me *Context) execUseTransaction(transactioner Transactioner, statements []*Statement) ([]sql.Result, error) {
	var saveResults []sql.Result

	for _, statement := range statements {
		result, err := transactioner.ExecStatement(statement)
		if err != nil {
			if rollbackError := transactioner.Rollback(); rollbackError != nil {
				return nil, rollbackError
			}
			return nil, err
		}
		saveResults = append(saveResults, result)
	}

	return saveResults, nil
}

//execWithoutTransaction execute statement without transaction
func (me *Context) execWithoutTransaction(execer Execer, statements []*Statement) ([]sql.Result, error) {
	var saveResults []sql.Result

	for _, statement := range statements {
		result, err := execer.ExecStatement(statement)
		if err != nil {
			return nil, err
		}
		saveResults = append(saveResults, result)
	}

	return saveResults, nil
}

//SaveChanges execute all defered statements to database
func (me *Context) SaveChanges() ([]sql.Result, error) {
	if me.ShouldUseTransaction() {
		if me.transaction == nil {
			newTransaction, err := me.BeginTransaction()
			if err != nil {
				return nil, err
			}

			results, err := me.execUseTransaction(newTransaction, me.Statements)
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

		results, err := me.execUseTransaction(me.transaction, me.Statements)
		me.ClearStatements()
		return results, err
	}

	results, err := me.execWithoutTransaction(me, me.Statements)
	me.ClearStatements()
	return results, err
}

//NewContext create new dbContext instance
func NewContext(dbClient *Client) *Context {
	return &Context{
		Client: *dbClient,
	}
}
