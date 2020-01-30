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

func Test_OrderService_CreateOrder(t *testing.T) {
	dbContext, err := initDBContext()
	defer dbContext.Close()

	if err != nil {
		t.Error(err.Error())
	}
	orderRepo := postgres.NewOrderRepository(dbContext)
	orderService := order.NewOrderService(orderRepo)
	newOrder := &order.Order{
		OrderNumber: "ORDER 01",
		OrderDate:   time.Now(),
		Total:       10000,
	}

	err = orderService.CreateOrder(context.Background(), newOrder)

	fmt.Print(newOrder.ID)
	if err != nil {
		t.Error(err.Error())
	}
}
