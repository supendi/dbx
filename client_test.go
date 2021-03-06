package dbx

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_Client_ExecStatement(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf("Fatal create db error: %s", err.Error())
		return
	}

	client := NewClient(db)

	var newID = uuid.New()
	//var newContext = context.Background()
	var sql = "INSERT INTO person (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)"

	statement := NewStatement(sql)
	statement.AddParameter("id", newID)
	statement.AddParameter("name", "Andi Setiawan")
	statement.AddParameter("created_at", time.Now())
	statement.AddParameter("updated_at", nil)

	result, err := client.ExecStatement(statement)
	if err != nil {
		t.Error("Error ", err.Error())
		return
	}

	if result == nil {
		t.Error("Error result nil")
		return
	}

	fmt.Println("Insert into person using client success")
}

func Test_Client_QueryStatement(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf("Fatal create db error: %s", err.Error())
		return
	}

	client := NewClient(db)

	var newID = uuid.New()
	//var newContext = context.Background()
	var insertSQL = "INSERT INTO person (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)"

	insertStatement := NewStatement(insertSQL)
	insertStatement.AddParameter("id", newID)
	insertStatement.AddParameter("name", "Dadang")
	insertStatement.AddParameter("created_at", time.Now())
	insertStatement.AddParameter("updated_at", nil)

	result, err := client.ExecStatement(insertStatement)

	if err != nil {
		t.Error("Error ", err.Error())
		return
	}

	if result == nil {
		t.Error("Error result nil")
		return
	}

	selectSQL := "SELECT * FROM person WHERE id=:id;"
	selectStatement := NewStatement(selectSQL)
	selectStatement.AddParameter("id", newID)

	rows, err := client.QueryStatement(selectStatement)

	if err != nil {
		t.Error("Query process error ", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		fmt.Println("Select using client succedd")
		return
	}
	t.Error("Query result must have records.")
}

func Test_Client_BeginTransaction(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf("Fatal create db error: %s", err.Error())
	}
	dbClient := NewClient(db)
	tx, err := dbClient.BeginTransaction()

	if err != nil {
		t.Errorf("Fatal begin transaction error: %s", err.Error())
	}

	if tx == nil {
		t.Errorf("Transaction expected to be not nil")
	}

	if dbClient.transaction == nil {
		t.Errorf("Fatal transaction error: transaction expected to be not <nil>")
	}
	dbClient.CompleteTransaction()
}

func Test_Client_GetTransaction(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf("Fatal create db error: %s", err.Error())
	}
	dbClient := NewClient(db)
	tx, err := dbClient.BeginTransaction()

	if err != nil {
		t.Errorf("Fatal begin transaction error: %s", err.Error())
	}

	if tx == nil {
		t.Errorf("Fatal transaction error: transaction expected to be not <nil>")
	}

	if dbClient.GetTransaction() == nil {
		t.Errorf("Fatal transaction error: transaction expected to be not <nil>")
	}

	fmt.Println("Test GetTransaction success")
}

func Test_Client_SetTransaction(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf("Fatal create db error: %s", err.Error())
	}
	dbClient := NewClient(db)
	tx, err := NewTransaction(db)

	if err != nil {
		t.Errorf("Fatal begin transaction error: %s", err.Error())
	}

	if tx == nil {
		t.Errorf("Fatal transaction error: transaction expected to be not <nil>")
	}

	dbClient.SetTransaction(tx)

	if dbClient.GetTransaction() == nil {
		t.Errorf("Fatal transaction error: transaction expected to be not <nil>")
	}

	fmt.Println("Test SetTransaction success")
}

func Test_Client_NewClient(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf("Fatal create db error: %s", err.Error())
	}
	dbClient := NewClient(db)
	if dbClient.DB == nil {
		t.Errorf("Fatal error: dbClient.DB expected to be not <nil>")
	}
	fmt.Println("Test NewClient success")
}
