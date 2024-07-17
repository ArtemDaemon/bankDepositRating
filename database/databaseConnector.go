package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func ExecQuery(query string) *sql.Rows {
	db := connect()
	defer db.Close()
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	return rows
}

func connect() *sql.DB {
	connStr := "user=postgres password=postgres dbname=bank_deposit_rating sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
