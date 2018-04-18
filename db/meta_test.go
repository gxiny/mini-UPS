package db

import (
	"testing"
)

func TestMeta(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	tx.Exec(`SAVEPOINT x`)
	worldId, err := GetMeta(tx, "world_id")
	if err == nil {
		t.Error("no error reading non-existent metadata")
	}
	tx.Exec(`ROLLBACK TO x`)

	err = SetMeta(tx, "world_id", "12345")
	if err != nil {
		t.Error(err)
	}
	worldId, err = GetMeta(tx, "world_id")
	if err != nil {
		t.Error(err)
	}
	if worldId != "12345" {
		t.Error("world_id != 12345")
	}
}
