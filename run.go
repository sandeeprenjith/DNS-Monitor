package main

import (
	"os"
	"fmt"
)

func handle_err(err error) {
        if err != nil {
                panic(err.Error())
        }
}

const usage = `Usage: ./dnsmon <path to config>
Run the script with config file(json) as argument.
`

func main() {
		CreateTable()
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
		servers, domains := Parse(file)
	for {
		for _, server := range servers {
			for _, domain := range domains {
				QueryAndSave(server, domain)
			}
		}
	}
}
