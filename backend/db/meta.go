package db

import (
	"database/sql"
)

var metaTable = sqlObject{
	`TABLE`, `meta`, `(
	key TEXT PRIMARY KEY,
	value TEXT
)`}

func SetMeta(tx *sql.Tx, key, value string) error {
	const sql = `INSERT INTO "meta" VALUES($1,$2) ` +
		`ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value`
	_, err := tx.Exec(sql, key, value)
	return err
}

func GetMeta(tx *sql.Tx, key string) (value string, err error) {
	const sql = `SELECT value FROM "meta" WHERE key = $1`
	err = tx.QueryRow(sql, key).Scan(&value)
	return
}
