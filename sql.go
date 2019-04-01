package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func CreateTable() {
	table := `CREATE TABLE IF NOT EXISTS "dnsdata" ( 
				"timestamp" TEXT,
				"server" TEXT,
				"qname" TEXT,
				"rtt" REAL,
				"rcode" TEXT);`
	db, err := sql.Open("sqlite3", "./db/dns.db")
	if err != nil {
		log.Fatal("[sql.go]", "Unable to open db 'dns.db'", err)
	}
	log.Print("[sql.go] Creating table 'dnsdata' if it doesn't already exist")
	stmt, err := db.Prepare(table)
	if err != nil {
		log.Print("[sql.go]", "Unable to prepare create table sql statement", err)
	}
	stmt.Exec()

	db.Close()
}

func InsertData(timestamp string, server string, qname string, rtt float64, rcode string) {
	sqlquery := `INSERT INTO "dnsdata" (
				"timestamp",
				"server",
				"qname",
				"rtt",
				"rcode") values (
				?, ?, ?, ?, ? );`
	db, err := sql.Open("sqlite3", "./db/dns.db")

	if err != nil {
		log.Fatal("[sql.go]", "Unable to open db 'dns.db'", err)
	}

	stmt, err := db.Prepare(sqlquery)
	if err != nil {
		log.Fatal("[sql.go]", "Unable prepare insert sql statement", err)
	}
	stmt.Exec(timestamp, server, qname, rtt, rcode)

}
