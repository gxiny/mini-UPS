package db

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	InitSchemas(tx)

	id, err := CreateUser(tx)
	if err != nil {
		t.Error(err)
	}

	_ = id
}
