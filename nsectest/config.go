package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
)

var verbose int
var resolvers []string
var basedomain string
var resolverfiles []string

type Resolver struct {
	name string
	ip   string
}
type Resolverlist []Resolver

var resolverlist Resolverlist

func getConfig() {
	// define and parse command line arguments
	pflag.CountVarP(&verbose, "verbose", "v", "print more information while running")
	pflag.StringSliceVarP(&resolvers, "resolver", "r", []string{}, "Resolver ip address and port ipv4:port [ipv6]:port")
	pflag.StringVarP(&basedomain, "domain", "d", "", "domain to test")
	pflag.StringSliceVarP(&resolverfiles, "resolvers", "R", []string{}, "file(s) with resolvers to test. Please see documentation.")
	pflag.Parse()

	if len(basedomain) == 0 {
		log.Fatal("Domain must be given")
	}
	if len(resolvers) == 0 && len(resolverfiles) == 0 {
		log.Fatal("Resolver or resolvers must be given")
	}

	if len(resolvers) > 0 {
		for _, resolver := range resolvers {
			resolverlist = append(resolverlist, Resolver{ip: resolver})
		}
	}
	if len(resolverfiles) > 0 {
		handleResolverfiles()
	}
}

func handleResolverfiles() {
	for _, filename := range resolverfiles {
		csv_file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		r := csv.NewReader(csv_file)
		_, _ = r.Read() // ignore first line
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			// record[0] name
			// record[1] ip address (ipv4 or ipv6)

			resolverlist = append(resolverlist, Resolver{name: record[0], ip: record[1]})
		}
		csv_file.Close()
	}
}
