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

	var pkg Package
	err = pkg.Create(tx, "abc123", Coord{3, 4}, 1)
	if err != nil {
		t.Error(err)
	}
}
