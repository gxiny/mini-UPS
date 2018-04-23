package db

import (
	"testing"
)

func TestPackage(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	var pkg Package
	err = pkg.Create(tx, []byte("abc123"), Coord{3, 4}, 1)
	if err != nil {
		t.Error(err)
	}

	err = pkg.SetDelivered(tx)
	if err != nil {
		t.Error(err)
	}
}
