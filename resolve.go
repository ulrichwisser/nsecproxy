package main

import (
	"log"
	"math/rand"
	"regexp"
	"strings"

	"github.com/miekg/dns"
)

func resolve(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	qname := req.Question[0].Name

	re := regexp.MustCompile("nsec.")
	switch re.FindString(dns.CanonicalName(qname)) {
	case "nsec0":
		log.Printf("routing NSEC0   request %s\n", qname)
		resolveNSEC0(config, w, req)
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
	case "nseca":
		log.Printf("routing NSECA   request %s\n", qname)
		resolveNSECA(config, w, req)
	case "nsecb":
		log.Printf("routing NSECB   request %s\n", qname)
		resolveNSECB(config, w, req)
	case "nsecc":
		log.Printf("routing NSECC   request %s\n", qname)
		resolveNSECC(config, w, req)
	case "nsecd":
		log.Printf("routing NSECD   request %s\n", qname)
		resolveNSECD(config, w, req)
	case "nsece":
		log.Printf("routing NSECE   request %s\n", qname)
		resolveNSECE(config, w, req)
	case "nsecf":
		log.Printf("routing NSECF   request %s\n", qname)
		resolveNSECF(config, w, req)
	case "nsecg":
		log.Printf("routing NSECG   request %s\n", qname)
		resolveNSECG(config, w, req)
	case "nsech":
		log.Printf("routing NSECH   request %s\n", qname)
		resolveNSECH(config, w, req)
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

// Returns a response from an upstream server
func getAnswer(upstream string, req *dns.Msg) *dns.Msg {
	c := &dns.Client{Net: "udp"}
	resp, _, err := c.Exchange(req, upstream)
	if err != nil {
		log.Println(err)
		return nil
	}
	return resp
}

// Replies to query with reponse from upstream
func resolveUpstream(upstream string, w dns.ResponseWriter, req *dns.Msg) {
	resp := getAnswer(upstream, req)
	if resp == nil {
		dns.HandleFailed(w, req)
	} else {
		w.WriteMsg(resp)
	}
	return
}

// strips all nsec and respective rrsig records from response
func stripDNSSEC(rrlist []dns.RR) []dns.RR {
	stripped := make([]dns.RR, 0)
	for _, rr := range rrlist {
		if rr.Header().Rrtype == dns.TypeNSEC ||
			rr.Header().Rrtype == dns.TypeRRSIG {
			continue
		}
		stripped = append(stripped, rr)
	}
	return stripped
}

// merge two rrlists, but filter out SOA from second list
func mergeRr(rrlist1, rrlist2 []dns.RR) []dns.RR {
	merged := make([]dns.RR, 0)
	for _, rr := range rrlist1 {
		merged = append(merged, rr)
	}
	for _, rr := range rrlist2 {
		if rr.Header().Rrtype == dns.TypeSOA {
			continue
		}
		if rr.Header().Rrtype == dns.TypeRRSIG && rr.(*dns.RRSIG).TypeCovered == dns.TypeSOA {
			continue
		}
		merged = append(merged, rr)
	}
	return merged
}

func resolveNSEC0(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC, req)
	resp1.Answer = stripDNSSEC(resp1.Answer)
	resp1.Ns = stripDNSSEC(resp1.Ns)
	resp1.Extra = stripDNSSEC(resp1.Extra)
	w.WriteMsg(resp1)
}

// response contains nsec and nsec records (in that order)
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
	w.WriteMsg(resp)
}

// response contains nsec3 and nsec records (in that order)
func resolveNSEC5(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC, req)
	resp3 := getAnswer(config.UpstreamNSEC3, req)

	resp := resp3.Copy()
	resp.Ns = mergeRr(resp1.Ns, resp3.Ns)
	w.WriteMsg(resp)
}

//
func resolveNSECA(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC, req)

	auth := make([]dns.RR, 0)
	for _, rr := range resp1.Ns {
		if rr.Header().Rrtype == dns.TypeNSEC ||
			(rr.Header().Rrtype == dns.TypeRRSIG && rr.(*dns.RRSIG).TypeCovered == dns.TypeNSEC) {
			continue
		}
		auth = append(auth, rr)
	}
	resp1.Ns = auth
	w.WriteMsg(resp1)
}

// NSEC does not cover label
func resolveNSECB(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC, req)
	req2 := req.Copy()
	labels := dns.SplitDomainName(req2.Question[0].Name)
	labels[0] = "d"
	req2.Question[0].Name = dns.Fqdn(strings.Join(labels, "."))
	resp2 := getAnswer(config.UpstreamNSEC, req2)
	resp1.Ns = resp2.Ns
	w.WriteMsg(resp1)
}

// NSEC3 does not cover label
func resolveNSECC(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC3, req)
	req2 := req.Copy()
	labels := dns.SplitDomainName(req2.Question[0].Name)
	labels[0] = "d"
	req2.Question[0].Name = dns.Fqdn(strings.Join(labels, "."))
	resp2 := getAnswer(config.UpstreamNSEC3, req2)
	resp1.Ns = resp2.Ns
	w.WriteMsg(resp1)
}

// NSEC3 does not cover label
func resolveNSECD(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC3, req)
	req2 := req.Copy()
	labels := dns.SplitDomainName(req2.Question[0].Name)
	labels[0] = "d"
	req2.Question[0].Name = dns.Fqdn(strings.Join(labels, "."))
	resp2 := getAnswer(config.UpstreamNSEC3, req2)
	resp1.Ns = resp2.Ns
	w.WriteMsg(resp1)
}

// NSEC does cover label, NSEC3 does not cover label
func resolveNSECE(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC, req)
	req2 := req.Copy()
	labels := dns.SplitDomainName(req2.Question[0].Name)
	labels[0] = "d"
	req2.Question[0].Name = dns.Fqdn(strings.Join(labels, "."))
	resp2 := getAnswer(config.UpstreamNSEC3, req2)
	resp1.Ns = mergeRr(resp1.Ns, resp2.Ns)
	w.WriteMsg(resp1)
}

// NSEC does cover label, NSEC3 does not cover label
func resolveNSECF(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC3, req)
	req2 := req.Copy()
	labels := dns.SplitDomainName(req2.Question[0].Name)
	labels[0] = "d"
	req2.Question[0].Name = dns.Fqdn(strings.Join(labels, "."))
	resp2 := getAnswer(config.UpstreamNSEC, req2)
	resp1.Ns = mergeRr(resp1.Ns, resp2.Ns)
	w.WriteMsg(resp1)
}

// data with nsec
func resolveNSECG(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC3, req)
	resp2 := getAnswer(config.UpstreamNSEC, req)
	resp1.Ns = resp2.Ns
	w.WriteMsg(resp1)
}

// data with nsec3
func resolveNSECH(config *Configuration, w dns.ResponseWriter, req *dns.Msg) {
	resp1 := getAnswer(config.UpstreamNSEC, req)
	resp2 := getAnswer(config.UpstreamNSEC3, req)
	resp1.Ns = resp2.Ns
	w.WriteMsg(resp1)
}
