package examples

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type OrderRepositoryer interface {
	GetAll(ctx context.Context) ([]*OrderModel, error)
	GetByID(ctx context.Context, orderID string) (*OrderModel, error)
	Add(ctx context.Context, order *OrderModel)
	Update(ctx context.Context, order *OrderModel)
	Delete(ctx context.Context, orderID string)
	SaveChanges() ([]sql.Result, error)
}

type OrderRepository struct {
	dbContext *DBContext
}

func (me *OrderRepository) mapToOrderDomain(order *Order) *OrderModel {
	return &OrderModel{
		ID:          order.ID,
		OrderNumber: order.OrderNumber,
		OrderDate:   order.OrderDate,
		Total:       order.Total,
	}
}

func (me *OrderRepository) mapToOrderPersistence(order *OrderModel) *Order {
	return &Order{
		ID:          order.ID,
		OrderNumber: order.OrderNumber,
		OrderDate:   order.OrderDate,
		Total:       order.Total,
		CreatedAt:   time.Now(),
	}
}

func (me *OrderRepository) GetAll(ctx context.Context) ([]*OrderModel, error) {
	orderRecords, err := me.dbContext.Order.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	orders := []*OrderModel{}

	for _, orderRecord := range orderRecords {
		order := me.mapToOrderDomain(orderRecord)
		orders = append(orders, order)
	}
	return orders, nil
}

func (me *OrderRepository) GetByID(ctx context.Context, orderID string) (*OrderModel, error) {
	orderRecord, err := me.dbContext.Order.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order := me.mapToOrderDomain(orderRecord)

	return order, nil
}

func (me *OrderRepository) Add(ctx context.Context, order *OrderModel) {
	var newOrderID = uuid.New().String()
	order.ID = newOrderID
	var newRecord = me.mapToOrderPersistence(order)
	me.dbContext.Order.Add(ctx, newRecord)
}

func (me *OrderRepository) Update(ctx context.Context, order *OrderModel) {
	var newRecord = me.mapToOrderPersistence(order)
	me.dbContext.Order.Update(ctx, newRecord)
}

func (me *OrderRepository) Delete(ctx context.Context, orderID string) {
	me.dbContext.Order.Delete(ctx, orderID)
}

func (me *OrderRepository) SaveChanges() ([]sql.Result, error) {
	return me.dbContext.SaveChanges()
}

func NewOrderRepository(dbContext *DBContext) OrderRepositoryer {
	return &OrderRepository{
		dbContext: dbContext,
	}
}
