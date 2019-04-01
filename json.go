package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var dat map[string][]string

func Parse(file string) ([]string, []string) {
	byt, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("[json.go]", "Unable to read config file", err)
	}
	if err = json.Unmarshal(byt, &dat); err != nil {
		log.Fatal("[json.go]", "Unable to unmarshal json data", err)
	}
	servers := dat["nameservers"]
	domains := dat["domains"]
	return servers, domains
}
