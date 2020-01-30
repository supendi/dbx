package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/supendi/dbx/examples/entities"
	"github.com/supendi/dbx/examples/order"
)

//OrderRepository implements order.OrderRepository
type OrderRepository struct {
	dbContext *entities.DBContext
}

//GetAll returns all order records
func (me *OrderRepository) GetAll(ctx context.Context) ([]*order.Order, error) {
	orderRecords, err := me.dbContext.Order.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	orders := []*order.Order{}
	for _, orderRecord := range orderRecords {
		orders = append(orders, &order.Order{
			ID:          orderRecord.ID,
			OrderNumber: orderRecord.OrderNumber,
			OrderDate:   orderRecord.OrderDate,
			Total:       orderRecord.Total,
			CreatedAt:   orderRecord.CreatedAt,
			UpdatedAt:   orderRecord.UpdatedAt,
		})
	}
	return orders, nil
}

//GetByID return single order record by order ID
func (me *OrderRepository) GetByID(ctx context.Context, orderID string) (*order.Order, error) {
	orderRecord, err := me.dbContext.Order.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	var newOrder *order.Order
	if orderRecord != nil {
		newOrder = &order.Order{
			ID:          orderRecord.ID,
			OrderNumber: orderRecord.OrderNumber,
			OrderDate:   orderRecord.OrderDate,
			Total:       orderRecord.Total,
			CreatedAt:   orderRecord.CreatedAt,
			UpdatedAt:   orderRecord.UpdatedAt,
		}
	}

	return newOrder, nil
}

//Add adds a new order into database
func (me *OrderRepository) Add(ctx context.Context, order *order.Order) (*order.Order, error) {
	order.ID = uuid.New().String()
	me.dbContext.Order.Add(&entities.Order{
		ID:          order.ID,
		OrderNumber: order.OrderNumber,
		OrderDate:   order.OrderDate,
		Total:       order.Total,
		CreatedAt:   order.CreatedAt,
	})
	_, err := me.dbContext.SaveChanges(ctx)
	return order, err
}

//Update updates existing order in database
func (me *OrderRepository) Update(ctx context.Context, order *order.Order) (*order.Order, error) {
	me.dbContext.Order.Update(&entities.Order{
		ID:          order.ID,
		OrderNumber: order.OrderNumber,
		OrderDate:   order.OrderDate,
		Total:       order.Total,
		CreatedAt:   order.CreatedAt,
	})
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
