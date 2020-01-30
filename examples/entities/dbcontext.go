package entities

import (
	"github.com/supendi/dbx"
)

//DBContext represent database context
type DBContext struct {
	*dbx.Context
	Order *OrderRepository
}

//NewDBContext returns new dbcontext
func NewDBContext(dbContext *dbx.Context) *DBContext {
	return &DBContext{
		Context: dbContext,
		Order:   NewOrderRepository(dbContext),
	}
}
