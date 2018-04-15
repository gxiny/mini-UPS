package db

import (
	"database/sql"
	"testing"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=passw0rd dbname=test")
	if err != nil {
		panic(err)
	}
	m.Run()
}
