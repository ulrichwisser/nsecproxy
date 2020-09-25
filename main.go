package main

import (
	"log"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

func main() {
	log.SetFlags(0)
	config := getConfig()

	// Start resolving
	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		go resolve(config, w, req)
	})

	//
	// Starting Servers
	//
	var wg sync.WaitGroup
	for _, ip := range config.IPlist {
		server := ip
		if strings.ContainsAny(":", server) {
			// IPv6 address
			server = "[" + server + "]:53"
		} else {
			// IPv4 address
			server = server + ":53"
		}

		log.Printf("starting to listen on %s", server)

		udpServer := &dns.Server{Addr: server, Net: "udp"}
		tcpServer := &dns.Server{Addr: server, Net: "tcp"}

		wg.Add(2)
		go func() {
			defer wg.Done()
			log.Fatal(udpServer.ListenAndServe())
		}()
		go func() {
			defer wg.Done()
			log.Fatal(tcpServer.ListenAndServe())
		}()
	}
	wg.Wait()

}
