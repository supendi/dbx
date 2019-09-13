package examples

import (
	"github.com/jmoiron/sqlx"
	"github.com/supendi/dbx"
)

type DBContext struct {
	*dbx.Context
	Order *OrderStatement
}

func NewDBContext(db *sqlx.DB) *DBContext {
	dbClient := dbx.NewClient(db)
	dbContext := dbx.NewContext(dbClient)
	return &DBContext{
		Context: dbContext,
		Order:   NewOrderStatement(dbContext),
	}
}
