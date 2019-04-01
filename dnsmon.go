package main

import (
	"github.com/miekg/dns"
	"log"
	"strconv"
	"strings"
)

func rcode_to_string(rc int) string {
	rcodes := map[int]string{
		0: "NOERROR",
		1: "FORMERR",
		2: "SERVFAIL",
		3: "NXDOMAIN",
		4: "NOTIMP",
		5: "REFUSED",
		6: "YXDOMAIN",
		7: "YXRRSET",
		8: "NOTAUTH",
		9: "NOTZONE"}
	return rcodes[rc]
}

func QueryAndSave(timestamp string, server string, qname string) {
	question := new(dns.Msg)
	question.SetQuestion(dns.Fqdn(qname), dns.TypeA)
	c := new(dns.Client)
	c.SingleInflight = true
	s_server := server + ":53"
	ans, t, err := c.Exchange(question, s_server)
	if err != nil {
		log.Print("[dnsmon.go]", "Unable to resolve DNS query", err)
	}
	rtt, err := strconv.ParseFloat(strings.Replace(t.String(), "ms", "", 1), 2)
	if err != nil {
		log.Print("[dnsmon.go]", err)
	}
	rcode := rcode_to_string(ans.Rcode)
	InsertData(timestamp, server, qname, rtt, rcode)
}
