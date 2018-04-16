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
	err = WithTx(db, func(tx *sql.Tx) error {
		return InitSchema(tx)
	})
	if err != nil {
		panic(err)
	}
	m.Run()
	err = WithTx(db, func(tx *sql.Tx) error {
		return DestroySchema(tx)
	})
	if err != nil {
		panic(err)
	}
}
