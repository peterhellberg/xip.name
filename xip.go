// Copyright 2014-2016 Peter Hellberg.
// Released under the terms of the MIT license.

// xip.name is a small name server which sends back any IP address found in the provided hostname.
//
// When queried for type A, it sends back the parsed IPv4 address.
// In the additional section the client host:port and transport are shown.
//
// Basic use pattern:
//
//    dig @xip.name foo.10.0.0.82.xip.name A
//
//    ; <<>> DiG 9.8.3-P1 <<>> @xip.name foo.10.0.0.82.xip.name A
//    ; (1 server found)
//    ;; global options: +cmd
//    ;; Got answer:
//    ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 13574
//    ;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
//    ;; WARNING: recursion requested but not available
//
//    ;; QUESTION SECTION:
//    ;foo.10.0.0.82.xip.name.		IN	A
//
//    ;; ANSWER SECTION:
//    foo.10.0.0.82.xip.name.	0	IN	A	10.0.0.82
//
//    ;; ADDITIONAL SECTION:
//    xip.name.		0	IN	TXT	"Client: 188.126.74.76:52575 (udp)"
//
//    ;; Query time: 27 msec
//    ;; SERVER: 188.166.43.179#53(188.166.43.179)
//    ;; WHEN: Wed Dec 31 02:55:51 2014
//    ;; MSG SIZE  rcvd: 128
//
// Initially based on the reflect example found at https://github.com/miekg/exdns
//
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/miekg/dns"
)

var (
	verbose = flag.Bool("v", false, "Verbose")
	fqdn    = flag.String("fqdn", "xip.name.", "FQDN to handle")
	addr    = flag.String("addr", ":53", "The addr to bind on")
	ip      = flag.String("ip", "188.166.43.179", "The IP of xip.name")

	ipPattern = regexp.MustCompile(`(\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	defaultIP net.IP
)

func main() {
	flag.Parse()

	defaultIP = net.ParseIP(*ip).To4()

	// Ensure that a FQDN is passed in (often the trailing . is omitted)
	*fqdn = dns.Fqdn(*fqdn)

	dns.HandleFunc(*fqdn, handleDNS)

	go serve(*addr, "tcp")
	go serve(*addr, "udp")

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%v) received, stopping\n", s)
}

func handleDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := &dns.Msg{}
	m.SetReply(r)

	if len(r.Question) == 0 {
		return
	}

	q := r.Question[0]
	t := dnsTXT(clientString(w.RemoteAddr()))
	rr := dnsRR(q.Name)

	switch q.Qtype {
	case dns.TypeTXT:
		m.Answer = append(m.Answer, t)
		m.Extra = append(m.Extra, rr)
	default:
		fallthrough
	case dns.TypeAAAA, dns.TypeA:
		m.Answer = append(m.Answer, rr)
		m.Extra = append(m.Extra, t)

	case dns.TypeAXFR, dns.TypeIXFR:
		c := make(chan *dns.Envelope)
		tr := new(dns.Transfer)
		defer close(c)

		err := tr.Out(w, r, c)
		if err != nil {
			if *verbose {
				fmt.Printf("%v\n", err)
			}

			return
		}

		soa := &dns.SOA{
			Hdr: dns.RR_Header{
				Name:   *fqdn,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    1440,
			},
			Ns:      *fqdn,
			Serial:  2014123101,
			Mbox:    *fqdn,
			Refresh: 21600,
			Retry:   7200,
			Expire:  604800,
			Minttl:  3600,
		}

		c <- &dns.Envelope{RR: []dns.RR{soa, t, rr, soa}}
		w.Hijack()

		return
	}

	if *verbose {
		fmt.Printf("%v\n", m.String())
	}

	w.WriteMsg(m)
}

func serve(addr, net string) {
	if err := newServer(addr, net).ListenAndServe(); err != nil {
		fmt.Printf("Failed to setup the %q server: %s\n", net, err.Error())
	}
}

func newServer(addr, net string) *dns.Server {
	return &dns.Server{
		Addr:       addr,
		Net:        net,
		TsigSecret: nil,
	}
}

func dnsRR(name string) (rr *dns.A) {
	rr = &dns.A{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    300,
		},
		A: defaultIP,
	}

	if ipStr := ipPattern.FindString(name); ipStr != "" {
		rr.A = net.ParseIP(ipStr).To4()
	}

	return rr
}

func dnsTXT(s string) *dns.TXT {
	return &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   *fqdn,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		Txt: []string{"Client: " + s},
	}
}

func clientString(a net.Addr) string {
	if ip, ok := a.(*net.UDPAddr); ok {
		return ip.String() + " (udp)"
	}

	if ip, ok := a.(*net.TCPAddr); ok {
		return ip.String() + " (tcp)"
	}

	return "unknown"
}
