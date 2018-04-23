package main

import (
	"bufio"
	"database/sql"
	"flag"
	"log"
	"net"

	_ "github.com/lib/pq"
	"gitlab.oit.duke.edu/rz78/ups/pb"
	"gitlab.oit.duke.edu/rz78/ups/pb/web"
)

var (
	dbOptions  = flag.String("db", "dbname=ups_server user=postgres password=passw0rd", "database options")
	listenAddr = flag.String("l", "8080", "listen address")
)

var database *sql.DB

func main() {
	var err error
	database, err = sql.Open("postgres", *dbOptions)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	req := new(web.Request)
	_, err := pb.ReadProto(bufio.NewReader(conn), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(conn.RemoteAddr(), req)
	resp := handleRequest(database, req)
	_, err = pb.WriteProto(conn, resp)
	if err != nil {
		log.Println(err)
	}
}
