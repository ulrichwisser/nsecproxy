package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/miekg/dns"
)

type rcodestats struct {
}

func makeRcode() {
	var rcodestats = make(map[string]map[string]int, 0)
	for m := range queries {
		rcodestats[m] = make(map[string]int, 0)
		for _, rcode := range getAllRcodes() {
			rcodestats[m][rcode] = 0
		}
	}

	for _, m := range config.data {
		for r := range m.Results {
			rstr := dns.RcodeToString[m.Results[r].Rcode]
			rcodestats[r][rstr]++
		}
	}

	// write data
	filename := filepath.Join(config.statsdir, "rcode.json")
	if config.verbose > 0 {
		log.Println("Writing rcode stats to ", filename)
	}
	b, err := json.Marshal(rcodestats)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	n, err := f.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	if n != len(b) {
		log.Fatalf("Marshaled data is %d bytes, written only %d bytes", len(b), n)
	}
	f.Sync()

}

func getAllRcodes() []string {
	var rcodes1 = make(map[string]int)

	for _, m := range config.data {
		for r := range m.Results {
			rstr := dns.RcodeToString[m.Results[r].Rcode]
			rcodes1[rstr] = 1
		}
	}

	rcodes2 := make([]string, len(rcodes1))

	i := 0
	for k := range rcodes1 {
		rcodes2[i] = k
		i++
	}

	return rcodes2
}
