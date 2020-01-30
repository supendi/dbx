package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/supendi/dbx"
	"github.com/supendi/dbx/examples/entities"
	"github.com/supendi/dbx/examples/order"
)

//OrderRepository implements order.OrderRepository
type OrderRepository struct {
	dbContext *entities.DBContext
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
		UpdatedAt:   order.UpdatedAt,
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

//Find return get list of orders filtered by specified filter
func (me *OrderRepository) Find(ctx context.Context, filter *order.OrderListFilter) ([]*order.Order, error) {
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	statement := dbx.NewStatement(`SELECT * FROM "order" WHERE order_number ILIKE :keyword ORDER BY created_at DESC LIMIT :limit`)
	statement.AddParameter("keyword", "%"+filter.Keyword+"%")
	statement.AddParameter("limit", filter.Limit)

	rows, err := me.dbContext.QueryStatementContext(ctx, statement)

	if err != nil {
		return nil, err
	}
	orders := []*order.Order{}
	for rows.Next() {
		order := &order.Order{}
		err = rows.Scan(&order.ID, &order.OrderNumber, &order.OrderDate, &order.Total, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

//NewOrderRepository returns new order repository instance
func NewOrderRepository(dbContext *entities.DBContext) *OrderRepository {
	return &OrderRepository{
		dbContext: dbContext,
	}
}
