// Package db contains database definitions of objects
// that are used in this project and related methods.
package db

import (
	"database/sql"
	"fmt"
)

type sqlObject struct {
	Typ  string
	Name string
	Def  string
}

func (s sqlObject) CreateSQL() string {
	return fmt.Sprintf(`CREATE %s "%s" %s`, s.Typ, s.Name, s.Def)
}

func (s sqlObject) DropSQL() string {
	return fmt.Sprintf(`DROP %s IF EXISTS "%s"`, s.Typ, s.Name)
}

var allSQL = [...]sqlObject{
	// types
	CoordType,
	// tables
	metaTable,
	UserTable,
	truckStatusType,
	TruckTable,
	PackageTable,
	PackageView,
	// indices
	TruckIdxWhId,
	TruckIdxStatus,
	PackageIdxUser,
	PackageIdxCtime,
	PackageIdxCtime1,
	PackageIdxWhId,
}

// InitSchema creates all objects in the database.
// It should be called when connecting to a new (empty)
// database.
func InitSchema(tx *sql.Tx) error {
	for _, obj := range allSQL {
		if _, err := tx.Exec(obj.CreateSQL()); err != nil {
			return err
		}
	}
	return nil
}

// DestroySchema deletes all objects in the database.
func DestroySchema(tx *sql.Tx) error {
	for i := len(allSQL) - 1; i >= 0; i-- {
		if _, err := tx.Exec(allSQL[i].DropSQL()); err != nil {
			return err
		}
	}
	return nil
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
