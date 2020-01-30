package order

import (
	"context"
	"time"
)

//Order represent order model
type Order struct {
	ID          string
	OrderNumber string
	OrderDate   time.Time
	Total       float64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

//OrderRepository bertugas untuk mengurus data access order
type OrderRepository interface {
	GetAll(ctx context.Context) ([]*Order, error)
	GetByID(ctx context.Context, orderID string) (*Order, error)
	Add(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, orderID string) error
}

//OrderService order business logic lay here
type OrderService struct {
	orderRepository OrderRepository
}

//CreateOrder creates a new order
func (me *OrderService) CreateOrder(ctx context.Context, order *Order) error {
	return me.orderRepository.Add(ctx, order)
}

//NewOrderService return a new order service instance
func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		orderRepository: repo,
	}
}
