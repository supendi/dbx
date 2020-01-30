package order_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/supendi/dbx"

	"github.com/supendi/dbx/examples/entities"
	"github.com/supendi/dbx/examples/order"
	"github.com/supendi/dbx/examples/order/postgres"
)

func getSqlxDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", "host=localhost dbname=dbx user=postgres password=irpan123 sslmode=disable")
	if err != nil {
		return db, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`TRUNCATE TABLE "order"`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDBContext() (*entities.DBContext, error) {
	db, err := getSqlxDb()
	if err != nil {
		return nil, err
	}

	client := dbx.NewClient(db)
	newDBcontext := dbx.NewContext(client)

	dbContext := entities.NewDBContext(newDBcontext)

	return dbContext, nil
}

func TestCreateOrder(t *testing.T) {
	dbContext, err := initDBContext()
	defer dbContext.Close()

	if err != nil {
		t.Error(err.Error())
	}
	orderRepo := postgres.NewOrderRepository(dbContext)
	orderService := order.NewOrderService(orderRepo)
	createRequest := &order.OrderCreateRequest{
		OrderNumber: "ORDER 01",
		OrderDate:   time.Now(),
		Total:       10000,
	}

	newOrder, err := orderService.CreateOrder(context.Background(), createRequest)
	fmt.Print(newOrder.ID)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestUpdateOrder(t *testing.T) {
	dbContext, err := initDBContext()

	dbContext.BeginTransaction()
	defer dbContext.Close()

	if err != nil {
		t.Error(err.Error())
	}
	orderRepo := postgres.NewOrderRepository(dbContext)
	orderService := order.NewOrderService(orderRepo)
	createRequest := &order.OrderCreateRequest{
		OrderNumber: "ORDER 01",
		OrderDate:   time.Now(),
		Total:       10000,
	}
	ctx := context.Background()
	createdOrder, err := orderService.CreateOrder(ctx, createRequest)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Print(createdOrder.ID + "\n")
	var newOrderNumber = "ORDER KE DUA"
	updateRequest := &order.OrderUpdateRequest{
		ID:          createdOrder.ID,
		OrderNumber: &newOrderNumber,
		OrderDate:   createdOrder.OrderDate,
		Total:       200,
	}

	updatedOrder, err := orderService.UpdateOrder(ctx, updateRequest)
	if err != nil {
		t.Error("Update error: " + err.Error())
	}
	err = dbContext.CompleteTransaction()
	if err != nil {
		t.Error("Transaction complete error: " + err.Error())
	}
	fmt.Print(updatedOrder)

}

func TestUpdateOrderMustError(t *testing.T) {
	dbContext, err := initDBContext()

	dbContext.BeginTransaction()
	defer dbContext.Close()

	if err != nil {
		t.Error(err.Error())
	}
	orderRepo := postgres.NewOrderRepository(dbContext)
	orderService := order.NewOrderService(orderRepo)
	createRequest := &order.OrderCreateRequest{
		OrderNumber: "ORDER 01",
		OrderDate:   time.Now(),
		Total:       10000,
	}
	ctx := context.Background()
	createdOrder, err := orderService.CreateOrder(ctx, createRequest)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Print(createdOrder.ID + "\n")
	//var newOrderNumber = "ORDER KE DUA"
	updateRequest := &order.OrderUpdateRequest{
		ID:          createdOrder.ID,
		OrderNumber: nil,
		OrderDate:   createdOrder.OrderDate,
		Total:       200,
	}

	updatedOrder, err := orderService.UpdateOrder(ctx, updateRequest)
	if err == nil {
		t.Error("Update process should return error but got nil")
	}
	err = dbContext.CompleteTransaction()
	if err != nil {
		t.Error("Transaction complete error: " + err.Error())
	}
	getRequest := &order.OrderGetRequest{
		ID: updatedOrder.ID,
	}
	fetchedOrder, err := orderService.GetOrder(ctx, getRequest)
	if err == nil {
		t.Error("Error should not be nil")
	}

	if fetchedOrder != nil {
		t.Error("Fecthed order should be nil. Because we want to make sure that  data wasn't stored on database")
		return
	}
}

func TestGetOrder(t *testing.T) {
	dbContext, err := initDBContext()

	defer dbContext.Close()

	if err != nil {
		t.Error(err.Error())
	}
	orderRepo := postgres.NewOrderRepository(dbContext)
	orderService := order.NewOrderService(orderRepo)
	createRequest := &order.OrderCreateRequest{
		OrderNumber: "ORDER 01",
		OrderDate:   time.Now(),
		Total:       10000,
	}
	ctx := context.Background()
	createdOrder, err := orderService.CreateOrder(ctx, createRequest)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Print(createdOrder.ID + "\n")
	getRequest := &order.OrderGetRequest{
		ID: createdOrder.ID,
	}

	fetchedOrder, err := orderService.GetOrder(ctx, getRequest)
	if err != nil {
		t.Error("Get error: " + err.Error())
	}

	if fetchedOrder == nil {
		t.Error("Fecthed order should not be nil")
		return
	}

	if fetchedOrder.ID != createdOrder.ID {
		t.Error("Fecthed order ID should be the same with created order ID")
	}
}

func TestListOrder(t *testing.T) {
	dbContext, err := initDBContext()

	defer dbContext.Close()

	if err != nil {
		t.Error(err.Error())
	}
	orderRepo := postgres.NewOrderRepository(dbContext)
	orderService := order.NewOrderService(orderRepo)
	createRequest := &order.OrderCreateRequest{
		OrderNumber: "ORDER 01",
		OrderDate:   time.Now(),
		Total:       10000,
	}
	ctx := context.Background()
	createdOrder, err := orderService.CreateOrder(ctx, createRequest)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Print(createdOrder.ID + "\n")
	listFilter := &order.OrderListFilter{
		Limit:   10,
		Keyword: "OR",
	}

	orders, err := orderService.ListOrder(ctx, listFilter)
	if err != nil {
		t.Error("list error: " + err.Error())
	}

	if orders == nil {
		t.Error("Fecthed order should not be nil")
		return
	}

	if len(orders) < 1 {
		t.Error("Orders length must be greater than 0")
	}

	listFilter.Keyword = "ODER" //it is a typo. It should be 'ORDER', which will be the cause of no data returned

	orders, err = orderService.ListOrder(ctx, listFilter)
	if err != nil {
		t.Error("list error: " + err.Error())
	}

	if orders == nil {
		t.Error("Fecthed order should not be nil")
		return
	}

	if len(orders) > 1 {
		t.Error("Orders length must be 0")
	}
}
