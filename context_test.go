package dbx

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_Context_AddStatement(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	dbClient := NewClient(db)
	dbContext := NewContext(dbClient)

	statement := NewStatement("SELECT * FROM person")
	dbContext.AddStatement(statement)

	if len(dbContext.Statements) != 1 {
		t.Errorf("Statements length must be one, but got %d", len(dbContext.Statements))
	}

	fmt.Println("AddStatement test succedd")
}

func Test_Context_ClearStatements(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	dbClient := NewClient(db)
	dbContext := NewContext(dbClient)
	statement := NewStatement("SELECT * FROM person")
	dbContext.AddStatement(statement)

	dbContext.ClearStatements()

	if len(dbContext.Statements) != 0 {
		t.Errorf("statements length must be 0, but got %d", len(dbContext.Statements))
	}

	fmt.Println("ClearStatements test succedd")
}

func Test_Context_ShouldUseTransaction1(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	dbClient := NewClient(db)
	dbContext := NewContext(dbClient)

	statement := NewStatement("SELECT * FROM person")

	dbContext.AddStatement(statement)
	dbContext.AddStatement(statement)

	if !dbContext.ShouldUseTransaction() {
		t.Errorf("Context should use transaction = true, but got = %v", dbContext.ShouldUseTransaction())
	}
	fmt.Println("ShouldUseTransaction test 1 succedd")
}

func Test_Context_ShouldUseTransaction2(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	dbClient := NewClient(db)
	dbContext := NewContext(dbClient)
	statement := NewStatement("SELECT * FROM person")

	dbContext.AddStatement(statement)
	tx, err := dbContext.BeginTransaction()
	if err != nil {
		t.Errorf(err.Error())
	}

	dbContext.Client.SetTransaction(tx)

	if !dbContext.ShouldUseTransaction() {
		t.Errorf("Context should use transaction = true even if its statements length is only 1, but got = %v", dbContext.ShouldUseTransaction())
	}

	fmt.Println("ShouldUseTransaction test 2 succeed")
}

func Test_Context_SaveChanges_SingleStatement(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	dbClient := NewClient(db)
	dbContext := NewContext(dbClient)

	var newID = uuid.New()
	//var newContext = context.Background()
	var insertSQL = "INSERT INTO person (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)"

	insertStatement := NewStatement(insertSQL)

	insertStatement.AddParameter("name", "Dadang")
	insertStatement.AddParameter("created_at", time.Now())
	insertStatement.AddParameter("updated_at", nil)
	insertStatement.AddParameter("id", newID) //prove that its ok if you add new parameter unsequential,

	dbContext.AddStatement(insertStatement)

	results, err := dbContext.SaveChanges(nil)

	if err != nil {
		t.Errorf(err.Error())
	}

	if results == nil {
		t.Errorf("Results expected not to be nil")
	}

	selectSQL := "SELECT * FROM person WHERE id=:id;"
	selectStatement := NewStatement(selectSQL)
	selectStatement.AddParameter("id", newID)

	rows, err := dbContext.QueryStatement(selectStatement)

	if err != nil {
		t.Error("Query process error ", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		fmt.Println("SaveChanges single statement test succeed")
		return
	}
	t.Error("Save changes failed, records not saved.")

}

func Test_Context_SaveChanges_MultiStatement_UseDefaultTransaction(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	dbClient := NewClient(db)
	dbContext := NewContext(dbClient)
	var newID1 = uuid.New()
	var newID2 = uuid.New()
	//var newContext = context.Background()
	var insertSQL = "INSERT INTO person (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)"

	insertStatement1 := NewStatement(insertSQL)
	insertStatement1.AddParameter("id", newID1)
	insertStatement1.AddParameter("name", "Dadang")
	insertStatement1.AddParameter("created_at", time.Now())
	insertStatement1.AddParameter("updated_at", nil)

	insertStatement2 := NewStatement(insertSQL)
	insertStatement2.AddParameter("id", newID2)
	insertStatement2.AddParameter("name", "Suhendra")
	insertStatement2.AddParameter("created_at", time.Now())
	insertStatement2.AddParameter("updated_at", nil)

	dbContext.AddStatement(insertStatement1)
	dbContext.AddStatement(insertStatement2)

	results, err := dbContext.SaveChanges(nil)

	if err != nil {
		t.Errorf(err.Error())
	}

	if results == nil {
		t.Errorf("Results expected not to be nil")
	}

	selectSQL := "SELECT * FROM person WHERE id=:id1 OR id =:id2;"
	selectStatement := NewStatement(selectSQL)
	selectStatement.AddParameter("id1", newID1)
	selectStatement.AddParameter("id2", newID1)

	rows, err := dbContext.QueryStatement(selectStatement)

	if err != nil {
		t.Error("Query process error ", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		fmt.Println("SaveChanges multi statement use default transaction test succeed")

		return
	}
	t.Error("Save changes failed, records not saved.")

}

func Test_Context_SaveChanges_MultiStatement_UseExternalTransaction(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	dbClient := NewClient(db)
	dbContext := NewContext(dbClient)
	tx, err := NewTransaction(db)

	if err != nil {
		t.Errorf(err.Error())
	}
	dbContext.SetTransaction(tx)

	var newID1 = uuid.New()
	var newID2 = uuid.New()
	//var newContext = context.Background()
	var insertSQL = "INSERT INTO person (id, name, created_at, updated_at) VALUES (:id, :name, :created_at, :updated_at)"

	insertStatement1 := NewStatement(insertSQL)
	insertStatement1.AddParameter("id", newID1)
	insertStatement1.AddParameter("name", "Dadang")
	insertStatement1.AddParameter("created_at", time.Now())
	insertStatement1.AddParameter("updated_at", nil)

	insertStatement2 := NewStatement(insertSQL)
	insertStatement2.AddParameter("id", newID2)
	insertStatement2.AddParameter("name", "Suhendra")
	insertStatement2.AddParameter("created_at", time.Now())
	insertStatement2.AddParameter("updated_at", nil)

	dbContext.AddStatement(insertStatement1)
	dbContext.AddStatement(insertStatement2)

	results, err := dbContext.SaveChanges(nil)
	tx.Commit()

	if err != nil {
		t.Errorf(err.Error())
	}

	if results == nil {
		t.Errorf("Results expected not to be nil")
	}

	selectSQL := "SELECT * FROM person WHERE id=:id1 OR id =:id2;"
	selectStatement := NewStatement(selectSQL)
	selectStatement.AddParameter("id1", newID1)
	selectStatement.AddParameter("id2", newID1)

	rows, err := dbContext.QueryStatement(selectStatement)

	if err != nil {
		t.Error("Query process error ", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		fmt.Println("SaveChanges multi statement use external transaction test succeed")
		return
	}
	t.Error("Save changes failed, records not saved.")
}

func Test_Context_SaveChanges_UseExternalTransaction_MustRolledBack(t *testing.T) {
	db, err := getSqlxDb()
	defer db.Close()
	if err != nil {
		t.Errorf(err.Error())
	}
	dbClient := NewClient(db)
	dbContext := NewContext(dbClient)
	tx, err := NewTransaction(db)

	if err != nil {
		t.Errorf(err.Error())
	}
	dbContext.SetTransaction(tx)

	var newID1 = uuid.New()
	//var newContext = context.Background()

	insertSQL := "INSERT INTO users (id, full_name, username, password, created_at, updated_at) VALUES (:id, :full_name, :username, :password, :created_at, :updated_at)"
	statement1 := NewStatement(insertSQL)
	statement1.AddParameter("id", newID1)
	statement1.AddParameter("full_name", "Andi Setiawan")
	statement1.AddParameter("username", "andi")
	statement1.AddParameter("password", "andi123")
	statement1.AddParameter("created_at", time.Now())
	statement1.AddParameter("updated_at", nil)

	statement2 := NewStatement(insertSQL)
	statement2.AddParameter("id", newID1) //duplicate id
	statement2.AddParameter("full_name", "Bowo")
	statement2.AddParameter("username", "bowo")
	statement2.AddParameter("password", "bowo123")
	statement2.AddParameter("created_at", time.Now())
	statement2.AddParameter("updated_at", nil)

	dbContext.AddStatement(statement1)
	dbContext.AddStatement(statement2)

	results, err := dbContext.SaveChanges(nil)

	if err == nil {
		t.Errorf("must return error.")
		if len(results) > 0 {
			t.Errorf("Results expected not have any values")
		}
	}

	selectSQL := "SELECT * FROM person WHERE id=:id1 OR id =:id2;"
	selectStatement := NewStatement(selectSQL)
	selectStatement.AddParameter("id1", newID1)
	selectStatement.AddParameter("id2", newID1)

	rows, err := dbContext.QueryStatement(selectStatement)

	if err != nil {
		t.Error("Query process error ", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		t.Error("SaveChanges multi statement use external transaction. Records must not be saved")
		return
	}
	fmt.Println("SaveChanges rollback test succeed")
}
