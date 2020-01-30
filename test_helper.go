package dbx

import (
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
	_, errEx := db.Exec("TRUNCATE TABLE person")
	if errEx != nil {
		return nil, err
	}
	return db, err
}
