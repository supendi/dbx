package examples

import (
	"context"
	"errors"
	"time"

	"github.com/supendi/dbx"
)

type Order struct {
	ID          string
	OrderNumber string
	OrderDate   time.Time
	Total       float64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type OrderStatement struct {
	dbContext *dbx.Context
}

func (me *OrderStatement) MapOrderToParam(statement *dbx.Statement, order *Order) {
	statement.AddParameter(`id`, order.ID)
	statement.AddParameter(`order_number`, order.OrderNumber)
	statement.AddParameter(`order_date`, order.OrderDate)
	statement.AddParameter(`total`, order.Total)
	statement.AddParameter(`created_at`, order.CreatedAt)
	statement.AddParameter(`updated_at`, order.UpdatedAt)
}

func (me *OrderStatement) GetAll(ctx context.Context) ([]*Order, error) {
	statement := dbx.NewStatement(`SELECT * FROM order`)
	me.dbContext.AddStatement(statement)

	rows, err := me.dbContext.QueryStatement(statement)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	orders := []*Order{}
	for rows.Next() {
		order := &Order{}
		err = rows.Scan(order.ID, order.OrderNumber, order.OrderDate, order.Total, order.CreatedAt, order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (me *OrderStatement) GetByID(ctx context.Context, orderID string) (*Order, error) {
	statement := dbx.NewStatement(`SELECT * FROM "order" WHERE id = :id`)
	statement.AddParameter(`id`, orderID)
	me.dbContext.AddStatement(statement)

	rows, err := me.dbContext.QueryStatement(statement)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		order := &Order{}
		err = rows.Scan(order.ID, order.OrderNumber, order.OrderDate, order.Total, order.CreatedAt, order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return order, nil
	}
	return nil, errors.New(`Record not found`)
}

func (me *OrderStatement) Add(ctx context.Context, order *Order) {
	statement := dbx.NewStatement(`INSERT INTO "order" (id, order_number, order_date, total, created_at, updated_at) VALUES (:id, :order_number, :order_date, :total, :created_at, :updated_at)`)
	me.MapOrderToParam(statement, order)

	me.dbContext.AddStatement(statement)
}

func (me *OrderStatement) Update(ctx context.Context, order *Order) {
	statement := dbx.NewStatement(`UPDATE "order" SET order_number = :order_number, order_date = :order_date, total = :total, created_at = :created_at, updated_at = :updated_at WHERE id=:id`)
	me.MapOrderToParam(statement, order)

	me.dbContext.AddStatement(statement)
}

func (me *OrderStatement) Delete(ctx context.Context, orderID string) {
	statement := dbx.NewStatement(`DELETE FROM "order" WHERE id = :id`)
	statement.AddParameter(`id`, orderID)

	me.dbContext.AddStatement(statement)
}

func NewOrderStatement(dbContext *dbx.Context) *OrderStatement {
	return &OrderStatement{
		dbContext: dbContext,
	}
}
