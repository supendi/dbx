package examples

import (
	"context"
	"time"
)

type OrderModel struct {
	ID          string
	OrderNumber string
	OrderDate   time.Time
	Total       float64
}

type OrderServicer interface {
	CreateOrder(ctx context.Context, order *OrderModel) error
}

type OrderService struct {
	orderRepository OrderRepositoryer
}

func (me *OrderService) CreateOrder(ctx context.Context, order *OrderModel) error {
	me.orderRepository.Add(ctx, order)
	_, err := me.orderRepository.SaveChanges()
	return err
}

func NewOrderService(repo OrderRepositoryer) OrderServicer {
	return &OrderService{
		orderRepository: repo,
	}
}
