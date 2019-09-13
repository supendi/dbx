package examples

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
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
	_, errEx := db.Exec(`TRUNCATE TABLE "order"`)
	if errEx != nil {
		return nil, err
	}
	return db, err
}

func initDBContext() (*DBContext, error) {
	db, err := getSqlxDb()

	if err != nil {
		return nil, err
	}

	dbContext := NewDBContext(db)
	return dbContext, nil
}

func Test_OrderService_CreateOrder(t *testing.T) {
	dbContext, err := initDBContext()

	if err != nil {
		t.Error(err.Error())
	}
	orderRepo := NewOrderRepository(dbContext)
	orderService := NewOrderService(orderRepo)
	newOrder := &OrderModel{
		OrderNumber: "ORDER 01",
		OrderDate:   time.Now(),
		Total:       10000,
	}

	err = orderService.CreateOrder(context.Background(), newOrder)
	if err != nil {
		t.Error(err.Error())
	}
}
