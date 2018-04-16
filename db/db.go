// Package db contains database schemas used in this project
// and helper functions.
package db

import (
	"database/sql"
)

// AllSQL contains SQL definitions of all database objects in this package.
var AllSQL = [...]string{
	CoordSQL,
	UserSQL,
	TruckSQL,
	PackageSQL,
}

// InitSchema creates all objects in the database.
// It should be called when connecting to a new (empty)
// database.
func InitSchema(tx *sql.Tx) error {
	for _, sql := range AllSQL {
		if _, err := tx.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

const destroySQL = `
DROP TABLE package;
DROP TABLE truck;
DROP TABLE "user";
DROP TYPE coordinate;
`

// DestroySchema deletes all objects in the database.
func DestroySchema(tx *sql.Tx) error {
	_, err := tx.Exec(destroySQL)
	return err
}

// WithTx encloses several database operations inside a transaction.
// If the operations succeed, the transaction is committed.
// Otherwise, it is rolled back.
func WithTx(db *sql.DB, query func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = query(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}
