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

var (
	dbOptions = flag.String("db", "dbname=ups_server user=postgres password=passw0rd", "database options")
	dbInit = flag.Bool("init", false, "init schema")
	dbDrop = flag.Bool("drop", false, "drop schema")
)

func main() {
	flag.Parse()
	database, err := sql.Open("postgres", *dbOptions)
	if err != nil {
		log.Fatal(err)
	}
	if *dbDrop {
		err = db.WithTx(database, func(tx *sql.Tx) error {
			return db.DestroySchema(tx)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	if *dbInit {
		err = db.WithTx(database, func(tx *sql.Tx) error {
			return db.InitSchema(tx)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
