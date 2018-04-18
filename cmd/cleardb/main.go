// Command cleardb empties the database created by
// the UPS server.
package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/lib/pq"
	"gitlab.oit.duke.edu/rz78/ups/db"
)

var dbOptions = flag.String("db", "dbname=test user=postgres password=passw0rd", "database options")

func main() {
	flag.Parse()
	database, err := sql.Open("postgres", *dbOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = db.WithTx(database, func(tx *sql.Tx) error {
		return db.DestroySchema(tx)
	})
	if err != nil {
		log.Fatal(err)
	}
}
