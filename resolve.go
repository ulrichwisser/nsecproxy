package main

import (
	"log"
	"math/rand"
	"regexp"

	"github.com/miekg/dns"
)

func resolve(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	qname := req.Question[0].Name

	re := regexp.MustCompile("nsec\\d")
	switch re.FindString(dns.CanonicalName(qname)) {
	case "nsec1":
		log.Printf("routing NSEC    request %s to       %s\n", qname, config.UpstreamNSEC)
		resolveUpstream(config.UpstreamNSEC, w, req)
	case "nsec3":
		log.Printf("routing NSEC3   request %s to       %s\n", qname, config.UpstreamNSEC3)
		resolveUpstream(config.UpstreamNSEC3, w, req)
	case "nsec4":
		log.Printf("routing NSEC4   request %s\n", qname)
		resolveNSEC4(config, w, req)
	case "nsec5":
		log.Printf("routing NSEC5   request %s\n", qname)
		resolveNSEC5(config, w, req)
	default:
		if rand.Intn(2) < 1 {
			log.Printf("routing DEFAULT request %s to NSEC  %s\n", qname, config.UpstreamNSEC)
			resolveUpstream(config.UpstreamNSEC, w, req)
		} else {
			log.Printf("routing DEFAULT request %s to NSEC3 %s\n", qname, config.UpstreamNSEC3)
			resolveUpstream(config.UpstreamNSEC3, w, req)
		}
	}
	return
}

func resolveUpstream(upstream string, w dns.ResponseWriter, req *dns.Msg) {
	resp := getAnswer(upstream, req)
	if resp == nil {
		dns.HandleFailed(w, req)
	} else {
		w.WriteMsg(resp)
	}
	return
}

func resolveNSEC4(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC, req)
	resp3 := getAnswer(config.UpstreamNSEC3, req)

	resp := resp1.Copy()
	for _, rr := range resp3.Ns {
		log.Printf("Adding to resp: %s\n", rr.String())
		if rr.Header().Rrtype == dns.TypeSOA {
			continue
		}
		if rr.Header().Rrtype == dns.TypeRRSIG && rr.(*dns.RRSIG).TypeCovered == dns.TypeSOA {
			continue
		}
		resp.Ns = append(resp.Ns, rr)
	}

	log.Println(resp.String())

	w.WriteMsg(resp)
}

func resolveNSEC5(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC, req)
	resp3 := getAnswer(config.UpstreamNSEC3, req)

	resp := resp3.Copy()
	for _, rr := range resp1.Ns {
		log.Printf("Adding to resp: %s\n", rr.String())
		if rr.Header().Rrtype == dns.TypeSOA {
			continue
		}
		if rr.Header().Rrtype == dns.TypeRRSIG && rr.(*dns.RRSIG).TypeCovered == dns.TypeSOA {
			continue
		}
		resp.Ns = append(resp.Ns, rr)
	}

	log.Println(resp.String())

	w.WriteMsg(resp)
}
func getAnswer(upstream string, req *dns.Msg) *dns.Msg {
	c := &dns.Client{Net: "udp"}
	resp, _, err := c.Exchange(req, upstream)
	if err != nil {
		log.Println(err)
		return nil
	}
	return resp
}
