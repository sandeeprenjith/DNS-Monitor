package main

/*
Here, we will be creating an API server. A user should be able to send an http
GET request and get a response in JSON format. The response should be as below.

{ "server": "1.1.1.1",
  "stats": [
			"domain": "google.com",
			"stats": {
					  "rcode": "NOERROR",
					  "rtt": 24.57},
			"domain": ".facebook.com",
			"stats": {
					  "rcode": "NOERROR",
					  "rtt": 23.1,},
			"domain": "youtube.com",
			"stats": {
					  "rcode": "NXDOMAIN",
					  "rtt": 24.57}
			]
}

The server should look at the sqlite database "dns.db" and fetch the latest
stats for the queried server. The db contains the DNS response statistics for
multiple domains queried against multiple servers.

The server should listen for the http GET request to the path /stats and expect
the server address in the "Body" of the GET request.

*/

import (
	// We will need the below packages for this.
	"database/sql"                  // to interact with the sql db
	_ "github.com/mattn/go-sqlite3" // This package makes it easier to interact
	// with sqlite3 db. But we don't need to
	// call the package. Hence the "_".
	"encoding/json" // to encode the response and to read the "Body" of the request.
	"log"           // Logging mainly for errors. logging format(errors) will be as below.
	"net/http"      // to make the http server.
	// [filename] Action which was being performed. Error: <error>
)

// Now let's create the structs. The first struct will be representing the
//individual stats for each domain queried.

type DomainStats struct {
	Rcode string  `json:"rcode"`
	RTT   float64 `json:"rtt"`
}

//Now let's create the struct for domain names and stats. stats for each domain
// will be of type "DomainStats".

type ServerStats struct {
	Domain string      `json:"domain"`
	Stats  DomainStats `json:"stats"`
}

// Now to create the struct that represents a server. Here each server can have
// multiple domains queried. Hence, the "Stats" will be a slice of the
// "ServerStats" type.

type Server struct {
	Server    string        `json:"server"`
	Timestamp string        `json:"timestamp"`
	Stats     []ServerStats `json:"stats"`
}

// Here we create a struct to hold the JSON Body from the API call

type Body struct {
	Server string `json:"server"`
}

// In the above structs the tags used (like `json:"server"` ) represent the
// keys when these structs are converted to JSON.

// Now we create a function to query the sqlite db and fetch the first results
// for a server. The function will accept the server address (string) as a
// parameter and return the resultset as JSON data as a slice of bytes.
func DBResults(server string) []byte {

	// Let's open the sqlite db "dns.db"
	db, err := sql.Open("sqlite3", "/db/dns.db")

	//In case of errors, log message and exit.
	if err != nil {
		log.Fatal("[api.go] Unable to open databse, 'dns.db'. Error:", err.Error())
	}

	// We need to close the db but let's defer that to the end of the function.
	defer db.Close()

	// Let's query the database and get the latest timestamp for the server.
	// First we prepare the statement.
	stmt := `select  timestamp from dnsdata where server ="` + server + `" order by timestamp DESC limit 1;`

	//Now we query
	rows, err := db.Query(stmt)

	// In case of errors log and exit.
	if err != nil {
		log.Fatal("[api.go] Unable to get data from SQL DB. Error: ", err.Error())
	}

	// Now we need to create a variable to store the timestamp from the db
	var tstamp string

	// Now we need to iterate through the data (here its only one, but the compiler
	//doesn't know that) and store the timestamp in the variable.
	for rows.Next() {
		rows.Scan(&tstamp)
	}

	//Lets close the "rows". Its best practice.
	rows.Close()

	// Now we query the results for the 'server' from the db based on the timestamp
	// we got from the last query.
	// First we prepare the statement.
	stmt = `select * from dnsdata where server =  "` + server + `" and timestamp = "` + tstamp + `";`

	// Now we query
	rows, err = db.Query(stmt)

	// You know the drill
	if err != nil {
		log.Fatal("[api.go] Unable to get data from SQL DB. Error: ", err.Error())
	}

	// Now we create more variables to hold the data
	var timestamp string
	var svr string // this is for the server
	var qname string
	var rtt float64
	var rcode string
	// Let's also create variables with the types of the structs we created
	var domainstats DomainStats
	var serverstats ServerStats
	var serverstatslice []ServerStats
	var Srvr Server // Since we already used the names "Server", "server" and "svr"

	// We iterate through the result and assign the data to the variables
	for rows.Next() {
		rows.Scan(&timestamp, &svr, &qname, &rtt, &rcode)
		// For each iteration(each row), we assign value to the struct variables.
		domainstats = DomainStats{Rcode: rcode, RTT: rtt}
		serverstats = ServerStats{Domain: qname, Stats: domainstats}
		// Now we append "serverstats" to the "serverstatslice" variable
		serverstatslice = append(serverstatslice, serverstats)
	}
	// Now we assign values to to the "Srvr" variable
	Srvr = Server{Server: svr, Stats: serverstatslice}

	//We convert the struct to JSON
	jsondata, err := json.MarshalIndent(Srvr, "", "    ")
	// Again with the errors
	if err != nil {
		log.Fatal("[api.go] Unable to convert data to json. Error: ", err.Error())
	}
	return jsondata
}

// Now lets create a function to handle the api calls.
// The function takes the http response writer and the http request as parameters
func serve(w http.ResponseWriter, r *http.Request) {

	// Lets define a variable for the server
	body := Body{}
	// Now we parse the json body from the API call and set the server
	err := json.NewDecoder(r.Body).Decode(&body)
	//As usual
	if err != nil {
		log.Print("[api.go] Unable to parse body from api call. Error:", err)
		// Set the HTTP header
		w.Header().Set("Content-Type", "application/json")
		// Write HTTP Status
		w.WriteHeader(http.StatusBadRequest)
		// Create response message
		err_response, _ := json.Marshal(`"error": "Unable to parse body. JSON format required"`)
		w.Write(err_response)
	} else {
		// Set the HTTP header
		w.Header().Set("Content-Type", "application/json")
		// Write HTTP Status
		w.WriteHeader(http.StatusOK)
		// Write HTTP Response
		w.Write(DBResults(body.Server))
	}
}

func main() {
	http.HandleFunc("/stats", serve)
	http.ListenAndServe(":8000", nil)
}
