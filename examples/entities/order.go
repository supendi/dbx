package entities

import (
	"context"
	"time"

	"github.com/supendi/dbx"
)

//Order represent order table
type Order struct {
	ID          string
	OrderNumber *string
	OrderDate   time.Time
	Total       float64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

//OrderRepository is
type OrderRepository struct {
	dbContext *dbx.Context
}

//setStatementParam sets statement parameters
func (me *OrderRepository) setStatementParam(statement *dbx.Statement, order *Order) {
	statement.AddParameter(`id`, order.ID)
	statement.AddParameter(`order_number`, order.OrderNumber)
	statement.AddParameter(`order_date`, order.OrderDate)
	statement.AddParameter(`total`, order.Total)
	statement.AddParameter(`created_at`, order.CreatedAt)
	statement.AddParameter(`updated_at`, order.UpdatedAt)
}

//GetAll gets all order records
func (me *OrderRepository) GetAll(ctx context.Context) ([]*Order, error) {
	statement := dbx.NewStatement(`SELECT * FROM order`)

	rows, err := me.dbContext.QueryStatementContext(ctx, statement)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	orders := []*Order{}
	for rows.Next() {
		order := &Order{}
		err = rows.Scan(&order.ID, &order.OrderNumber, &order.OrderDate, &order.Total, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

//GetByID gets order by order ID
func (me *OrderRepository) GetByID(ctx context.Context, orderID string) (*Order, error) {
	statement := dbx.NewStatement(`SELECT * FROM "order" WHERE id = :id`)
	statement.AddParameter(`id`, orderID)

	rows, err := me.dbContext.QueryStatementContext(ctx, statement)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		order := &Order{}
		err = rows.Scan(&order.ID, &order.OrderNumber, &order.OrderDate, &order.Total, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return order, nil
	}
	return nil, nil
}

//Add adds new order
func (me *OrderRepository) Add(order *Order) {
	statement := dbx.NewStatement(`INSERT INTO "order" (id, order_number, order_date, total, created_at, updated_at) VALUES (:id, :order_number, :order_date, :total, :created_at, :updated_at)`)
	me.setStatementParam(statement, order)

	me.dbContext.AddStatement(statement)
}

//Update updates existing order
func (me *OrderRepository) Update(order *Order) {
	statement := dbx.NewStatement(`UPDATE "order" SET order_number = :order_number, order_date = :order_date, total = :total, created_at = :created_at, updated_at = :updated_at WHERE id=:id`)
	me.setStatementParam(statement, order)

	me.dbContext.AddStatement(statement)
}

//Delete deletes existing order
func (me *OrderRepository) Delete(orderID string) {
	statement := dbx.NewStatement(`DELETE FROM "order" WHERE id = :id`)
	statement.AddParameter(`id`, orderID)

	me.dbContext.AddStatement(statement)
}

//NewOrderRepository create new order table instance
func NewOrderRepository(dbContext *dbx.Context) *OrderRepository {
	return &OrderRepository{
		dbContext: dbContext,
	}
}
