package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/supendi/dbx/examples/entities"
)

//OrderRepository implements order.OrderRepository
type OrderRepository struct {
	dbContext *entities.DBContext
}

//GetAll returns all order records
func (me *OrderRepository) GetAll(ctx context.Context) ([]*entities.Order, error) {
	orderRecords, err := me.dbContext.Order.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	orders := []*entities.Order{}

	for _, orderRecord := range orderRecords {
		orders = append(orders, orderRecord)
	}
	return orders, nil
}

//GetByID return single order record by order ID
func (me *OrderRepository) GetByID(ctx context.Context, orderID string) (*entities.Order, error) {
	orderRecord, err := me.dbContext.Order.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return orderRecord, nil
}

//Add adds new order into database
func (me *OrderRepository) Add(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	var newOrderID = uuid.New().String()
	order.ID = newOrderID
	me.dbContext.Order.Add(order)
	_, err := me.dbContext.SaveChanges(ctx)
	return order, err
}

//Update updates existing order in database
func (me *OrderRepository) Update(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	me.dbContext.Order.Update(order)
	_, err := me.dbContext.SaveChanges(ctx)
	return order, err
}

//Delete deletes existing order
func (me *OrderRepository) Delete(ctx context.Context, orderID string) error {
	me.dbContext.Order.Delete(orderID)
	_, err := me.dbContext.SaveChanges(ctx)
	return err
}

//NewOrderRepository returns new order repository instance
func NewOrderRepository(dbContext *entities.DBContext) *OrderRepository {
	return &OrderRepository{
		dbContext: dbContext,
	}
}
