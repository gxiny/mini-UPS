package db

import (
	"database/sql"
)

const UserSQL = `
CREATE TABLE "user" (
	id BIGSERIAL PRIMARY KEY
);`

// CreateUser returns the ID of a newly created user.
func CreateUser(tx *sql.Tx) (id int64, err error) {
	sql := `INSERT INTO "user" DEFAULT VALUES RETURNING id`
	err = tx.QueryRow(sql).Scan(&id)
	return
}
