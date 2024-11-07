package db_client

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DBClient *sql.DB

func InitializeDBConnection() {
	db, err := sql.Open("postgres", "postgres://postgres:test123@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	DBClient = db
}
