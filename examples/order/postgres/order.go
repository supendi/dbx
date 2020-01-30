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

//mapToOrderDomain return a new order domain instance from an entity object
func (me *OrderRepository) mapEntityToDomain(entity *entities.Order) *order.Order {
	return &order.Order{
		ID:          entity.ID,
		OrderNumber: entity.OrderNumber,
		OrderDate:   entity.OrderDate,
		Total:       entity.Total,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

//mapToOrderEntity return a new order entity instance from an domain model object
func (me *OrderRepository) mapDomainToEntity(ord *order.Order) *entities.Order {
	return &entities.Order{
		ID:          ord.ID,
		OrderNumber: ord.OrderNumber,
		OrderDate:   ord.OrderDate,
		Total:       ord.Total,
		CreatedAt:   ord.CreatedAt,
		UpdatedAt:   ord.UpdatedAt,
	}
}

//GetAll returns all order records
func (me *OrderRepository) GetAll(ctx context.Context) ([]*order.Order, error) {
	orderRecords, err := me.dbContext.Order.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	orders := []*order.Order{}

	for _, orderRecord := range orderRecords {
		order := me.mapEntityToDomain(orderRecord)
		orders = append(orders, order)
	}
	return orders, nil
}

//GetByID return single order record by order ID
func (me *OrderRepository) GetByID(ctx context.Context, orderID string) (*order.Order, error) {
	orderRecord, err := me.dbContext.Order.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	var order *order.Order
	if orderRecord != nil {
		order = me.mapEntityToDomain(orderRecord)
	}

	return order, nil
}

//Add adds new order into database
func (me *OrderRepository) Add(ctx context.Context, order *order.Order) (*order.Order, error) {
	var newOrderID = uuid.New().String()
	order.ID = newOrderID
	var newRecord = me.mapDomainToEntity(order)
	me.dbContext.Order.Add(newRecord)
	_, err := me.dbContext.SaveChanges(ctx)
	return order, err
}

//Update updates existing order in database
func (me *OrderRepository) Update(ctx context.Context, order *order.Order) (*order.Order, error) {
	var entity = me.mapDomainToEntity(order)
	me.dbContext.Order.Update(entity)
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
