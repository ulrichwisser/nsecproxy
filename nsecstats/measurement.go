package main

import (
	"github.com/miekg/dns"
)

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
