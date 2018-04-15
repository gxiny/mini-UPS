package db

import (
	"testing"
)

func TestCreatePackage(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	err = InitSchemas(tx)
	if err != nil {
		t.Fatal(err)
	}

	id, err := CreatePackage(tx, "abc123", Coord{3, 4}, 1)
	if err != nil {
		t.Error(err)
	}

	_ = id
}
