package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func dont_panic(err error) {
	if err != nil {
		return
	}
}
func CreateTable() {
	table := `CREATE TABLE "dnsdata" ( 
				"timestamp" TEXT,
				"server" TEXT,
				"qname" TEXT,
				"rtt" REAL,
				"rcode" TEXT);`
	query := `SELECT * FROM dnsdata;`
	db, err :=sql.Open("sqlite3", "./dns.db")
	handle_err(err)
	rows, err := db.Query(query)
	if rows != nil {
			return
	}
	if err != nil {
			stmt, err := db.Prepare(table)
			handle_err(err)
			stmt.Exec()
	}
	db.Close()
}

func InsertData(server string, qname string, rtt float64, rcode string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	sqlquery := `INSERT INTO "dnsdata" (
				"timestamp",
				"server",
				"qname",
				"rtt",
				"rcode") values (
				?, ?, ?, ?, ? );` 
	db, err :=sql.Open("sqlite3", "./dns.db")
	handle_err(err)
	stmt, err := db.Prepare(sqlquery)
	handle_err(err)
	stmt.Exec(timestamp, server, qname, rtt, rcode)

}

