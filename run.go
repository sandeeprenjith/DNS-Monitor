package main

import (
	"fmt"
	"os"
	"time"
)

const usage = `Usage: ./dnsmon <path to config>
Run the script with config file(json) as argument.`

func QandS(servers []string, domains []string) {
	for {
		for _, server := range servers {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			for _, domain := range domains {
				QueryAndSave(timestamp, server, domain) //see dnsmon.go
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func main() {
	CreateTable() //see sql.go
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	} else if os.Args[1] == "-h" {
		fmt.Println(usage)
		return
	} else if os.Args[1] == "--help" {
		fmt.Println(usage)
		return
	}
	file := os.Args[1]
	servers, domains := Parse(file) //see json.go
	QandS(servers, domains)
}
