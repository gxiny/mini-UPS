package db

import (
	"testing"
)

func TestUserCreate(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	var user User
	err = user.Create(tx)
	if err != nil {
		t.Error(err)
	}
}
