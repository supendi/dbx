package order

import (
	"context"
	"errors"
	"time"
)

//ErrOrderNotFound is an error if order is not found on data storage
var ErrOrderNotFound = errors.New("Order is not found")

//Order represent order model
type Order struct {
	ID          string
	OrderNumber *string
	OrderDate   time.Time
	Total       float64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

//OrderCreateRequest represent the request model for creating an order
type OrderCreateRequest struct {
	OrderNumber string
	OrderDate   time.Time
	Total       float64
	CreatedAt   time.Time
}

//OrderUpdateRequest represent the request model for updating an existing order
type OrderUpdateRequest struct {
	ID          string
	OrderNumber *string
	OrderDate   time.Time
	Total       float64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

//OrderRepository bertugas untuk mengurus data access order
type OrderRepository interface {
	GetAll(ctx context.Context) ([]*Order, error)
	GetByID(ctx context.Context, orderID string) (*Order, error)
	Add(ctx context.Context, order *Order) (*Order, error)
	Update(ctx context.Context, order *Order) (*Order, error)
	Delete(ctx context.Context, orderID string) error
}

//OrderService order business logic lay here
type OrderService struct {
	orderRepository OrderRepository
}

//CreateOrder creates a new order
func (me *OrderService) CreateOrder(ctx context.Context, request *OrderCreateRequest) (*Order, error) {
	newOrder := &Order{
		OrderNumber: &request.OrderNumber,
		OrderDate:   request.OrderDate,
		Total:       request.Total,
		CreatedAt:   time.Now(),
	}
	return me.orderRepository.Add(ctx, newOrder)
}

//UpdateOrder updates an existing order
func (me *OrderService) UpdateOrder(ctx context.Context, request *OrderUpdateRequest) (*Order, error) {
	existingOrder, err := me.orderRepository.GetByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	if existingOrder == nil {
		return nil, ErrOrderNotFound
	}

	updatedDate := &time.Time{}
	*updatedDate = time.Now()

	existingOrder.OrderNumber = request.OrderNumber
	existingOrder.OrderDate = request.OrderDate
	existingOrder.Total = request.Total
	existingOrder.UpdatedAt = updatedDate

	return me.orderRepository.Update(ctx, existingOrder)
}

//NewOrderService return a new order service instance
func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		orderRepository: repo,
	}
}
