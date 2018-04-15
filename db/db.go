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

// InitSchemas creates all schemas in the database.
// It should be called when connecting to a new (empty)
// database.
func InitSchemas(tx *sql.Tx) error {
	for _, sql := range AllSQL {
		if _, err := tx.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

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
