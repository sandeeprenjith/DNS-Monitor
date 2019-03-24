package main

import (
	"encoding/json"
	"io/ioutil"
)

var dat map[string][]string

func Parse(file string) ([]string, []string){
	byt, err := ioutil.ReadFile(file)
	if err != nil{
		panic(err.Error())
	}
	if err = json.Unmarshal(byt, &dat); err != nil {
        panic(err)
    }
	servers := dat["nameservers"]
	domains:= dat["domains"]
	return servers, domains
}
