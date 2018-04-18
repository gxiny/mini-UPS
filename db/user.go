package db

import (
	"database/sql"
)

var UserTable = sqlObject{
	`TABLE`, `user`, `(
	id BIGSERIAL PRIMARY KEY
)`}

type User int64

// Create creates a new user.
// The receiver is modified to the ID of the new user.
func (id *User) Create(tx *sql.Tx) error {
	sql := `INSERT INTO "user" DEFAULT VALUES RETURNING id`
	return tx.QueryRow(sql).Scan(id)
}
