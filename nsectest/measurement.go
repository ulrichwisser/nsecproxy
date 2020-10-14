package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
)

const timeout = 30

type result struct {
	FlagQR     bool `json:"flagqr"`
	FlagRD     bool `json:"flagrd"`
	FlagRA     bool `json:"flagra"`
	FlagAD     bool `json:"flagad"`
	FlagDO     bool `json:"flagdo"`
	FlagAA     bool `json:"flagaa"`
	FlagTC     bool `json:"flagtc"`
	FlagCD     bool `json:"flagcd"`
	Response   bool `json:"response"`
	Rcode      int  `json:"rcode"`
	Answer     int  `json:"answer"`
	Authority  int  `json:"authority"`
	Additional int  `json:"additional"`
}

type results map[string]result

type query struct {
	qname string
	qtype uint16
}

var queries = map[string]query{
	"NoDataNoWildcardNsec1":                 {"nsec1.a", dns.TypeAAAA},
	"NoDataNoWildcardNsec3":                 {"nsec3.a", dns.TypeAAAA},
	"NoDataNoWildcardNsec4":                 {"nsec4.a", dns.TypeAAAA},
	"NoDataNoWildcardNsec5":                 {"nsec5.a", dns.TypeAAAA},
	"NoDataWildcardNsec1":                   {"nsec1.b", dns.TypeAAAA},
	"NoDataWildcardNsec3":                   {"nsec1.b", dns.TypeAAAA},
	"NoDataWildcardNsec4":                   {"nsec1.b", dns.TypeAAAA},
	"NoDataWildcardNsec5":                   {"nsec1.b", dns.TypeAAAA},
	"NameErrorNoWildcardNsec1":              {"nsec1.c", dns.TypeAAAA},
	"NameErrorNoWildcardNsec3":              {"nsec3.c", dns.TypeAAAA},
	"NameErrorNoWildcardNsec4":              {"nsec4.c", dns.TypeAAAA},
	"NameErrorNoWildcardNsec5":              {"nsec5.c", dns.TypeAAAA},
	"EmptyNonTerminalNsec1":                 {"nsec1.d", dns.TypeAAAA},
	"EmptyNonTerminalNsec3":                 {"nsec1.d", dns.TypeAAAA},
	"EmptyNonTerminalNsec4":                 {"nsec1.d", dns.TypeAAAA},
	"EmptyNonTerminalNsec5":                 {"nsec1.d", dns.TypeAAAA},
	"FailureNoNsecNoNsec3":                  {"a.nseca", dns.TypeAAAA},
	"FailureNsecDoesNotCoverLabel":          {"b.nsecb", dns.TypeAAAA},
	"FailureNsec3DoesNotCoverLabel":         {"b.nsecc", dns.TypeAAAA},
	"FailureNsecAndNsec3DoesNotCoverLabel":  {"b.nsecd", dns.TypeAAAA},
	"FailureNsecDoesCoverLabelNsec3doesNot": {"b.nsece", dns.TypeAAAA},
	"FailureNsecDoesNotCoverLabelNsec3Does": {"b.nsecf", dns.TypeAAAA},
	"FailureDataWithNsec":                   {"b.nsecg", dns.TypeTXT},
	"FailureDataWithNsec3":                  {"b.nsech", dns.TypeTXT},
}

type measurement struct {
	Name    string
	Results results
}

func measure(resolver Resolver) (measurement, error) {
	if verbose > 0 {
		log.Printf("Starting measurement on %s", resolver)
	}
	var measurement = measurement{Name: resolver.name, Results: make(results, 0)}

	for label, q := range queries {
		r, err := resolve(ip2dial(resolver.ip), q)
		if err != nil {
			return measurement, err
		}
		measurement.Results[label] = r
	}
	return measurement, nil
}

// resolv will send a query and return the result
func resolve(resolver string, m query) (result, error) {
	if verbose > 2 {
		fmt.Printf("Resolving %s %s (%d)\n", m.qname, dns.TypeToString[m.qtype], m.qtype)
	}

	// Setting up query
	query := new(dns.Msg)
	query.RecursionDesired = true
	query.Question = make([]dns.Question, 1)
	query.SetQuestion(dns.Fqdn(strings.Join([]string{m.qname, basedomain}, ".")), m.qtype)

	// Setting up resolver
	client := new(dns.Client)
	client.ReadTimeout = timeout * 1e9

	// make the query and wait for answer
	r, _, err := client.Exchange(query, resolver)

	// check for errors
	if err != nil {
		log.Printf("Error resolving %s %s (%d): %s\n", m.qname, dns.TypeToString[m.qtype], m.qtype, err)
		return result{}, err
	}
	if r == nil {
		log.Printf("No answer %s %s (%d)\n", m.qname, dns.TypeToString[m.qtype], m.qtype)
		return result{}, err
	}

	var result result
	result.Response = r.Response
	result.FlagRD = r.RecursionDesired
	result.FlagRA = r.RecursionAvailable
	result.FlagAD = r.AuthenticatedData
	result.FlagAA = r.Authoritative
	result.FlagTC = r.Truncated
	result.FlagCD = r.CheckingDisabled
	if r.IsEdns0() != nil {
		result.FlagDO = r.IsEdns0().Do()
	} else {
		result.FlagDO = false
	}
	result.Rcode = r.Rcode
	result.Answer = len(r.Answer)
	result.Authority = len(r.Ns)
	result.Additional = len(r.Extra)

	return result, nil
}

func ip2dial(ip string) (server string) {
	server = ip
	if strings.ContainsAny(":", server) {
		// IPv6 address
		server = "[" + server + "]:53"
	} else {
		server = server + ":53"
	}
	return
}
