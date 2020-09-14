package main

import (
	"log"

	"github.com/miekg/dns"
)

func main() {
	config := getConfig()

	//
	// Starting Servers
	//
	if config.Verbose > 0 {
		log.Println("Startung UDP server ", config.ListenUDP)
	}
	udpServer := &dns.Server{Addr: config.ListenUDP, Net: "udp"}

	if len(config.ListenTCP) > 0 {
		if config.Verbose > 0 {
			log.Println("Startung TCP server ", config.ListenTCP)
		}
		//tcpServer := &dns.Server{Addr: config.ListenTCP, Net: "tcp"}
	} else {
		if config.Verbose > 0 {
			log.Println("No TCP server ")
		}
	}
	// Start resolving
	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		go resolve(config, w, req)
	})

	//go func() {
	log.Fatal(udpServer.ListenAndServe())
	//}()
	//log.Fatal(tcpServer.ListenAndServe())

}
