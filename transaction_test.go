package dbx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_Transaction_ExecStatement(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf("Fatal create db error: %s", err.Error())
		return
	}

	dbTrans, err := NewTransaction(db)
	if err != nil {
		t.Errorf("Fatal db transaction error: %s", err.Error())
		return
	}

	var newId = uuid.New()
	var newContext = context.Background()
	var sql = "INSERT INTO unit_testx (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)"

	statement := NewStatement(newContext, sql)
	statement.AddParameter("id", newId)
	statement.AddParameter("name", "Irpan Supendi")
	statement.AddParameter("created_at", time.Now())
	statement.AddParameter("updated_at", nil)

	result, err := dbTrans.ExecStatement(statement)
	dbTrans.Commit()
	if err != nil {
		t.Error("Error ", err.Error())
		return
	}

	if result == nil {
		t.Error("Error result nil")
		return
	}

	fmt.Println("Insert into unit_testx using transaction success")
}

func Test_Transaction_QueryStatement(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf("Fatal create db error: %s", err.Error())
		return
	}

	dbTrans, err := NewTransaction(db)
	if err != nil {
		t.Errorf("Fatal db transaction error: %s", err.Error())
		return
	}

	var newId = uuid.New()
	var newContext = context.Background()
	var insertSql = "INSERT INTO unit_testx (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)"

	statement := NewStatement(newContext, insertSql)
	statement.AddParameter("id", newId)
	statement.AddParameter("name", "Irpan Supendi")
	statement.AddParameter("created_at", time.Now())
	statement.AddParameter("updated_at", nil)

	result, err := dbTrans.ExecStatement(statement)

	if err != nil {
		t.Error("Error ", err.Error())
		return
	}

	if result == nil {
		t.Error("Error result nil")
		return
	}

	err = dbTrans.CommitAndStartOver()
	if err != nil {
		t.Error("Commit error ", err.Error())
		return
	}

	selectSql := "SELECT * FROM unit_testx WHERE id=:id;"
	queryStatement := NewStatement(newContext, selectSql)
	queryStatement.AddParameter("id", newId)
	if dbTrans.IsComplete() {
		err = dbTrans.StartOver()
		if err != nil {
			t.Error("Start new transaction error ", err.Error())
			return
		}
	}

	rows, err := dbTrans.QueryStatement(queryStatement)

	if err != nil {
		t.Error("Query process error ", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		fmt.Println("Select using transaction succed")
		return
	}
	t.Error("Query result must have records.")
}
