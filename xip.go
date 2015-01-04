// Copyright 2014-2015 Peter Hellberg.
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
	port    = flag.String("p", "53", "The port to bind on")
	ip      = flag.String("ip", "188.166.43.179", "The IP of xip.name")

	ipPattern = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	defaultIP net.IP
)

func main() {
	flag.Parse()

	defaultIP = net.ParseIP(*ip).To4()

	dns.HandleFunc(*fqdn, handleDNS)

	go serve("tcp")
	go serve("udp")

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case s := <-sig:
			fmt.Printf("\nSignal (%d) received, stopping\n", s)
			break loop
		}
	}
}

func handleDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
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
			return
		}

		soa, _ := dns.NewRR(`xip.name. xip.name. 0 IN SOA 2014123101 21600 7200 604800 3600`)

		c <- &dns.Envelope{RR: []dns.RR{soa, t, rr, soa}}
		w.Hijack()

		return
	}

	if *verbose {
		fmt.Printf("%v\n", m.String())
	}

	w.WriteMsg(m)
}

func serve(net string) {
	err := newServer(net).ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())
	}
}

func newServer(net string) *dns.Server {
	return &dns.Server{
		Addr:       ":" + *port,
		Net:        net,
		TsigSecret: nil,
	}
}

func dnsRR(name string) (rr dns.RR) {
	rr = new(dns.A)
	rr.(*dns.A).Hdr = dns.RR_Header{
		Name:   name,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    300,
	}

	if ipStr := ipPattern.FindString(name); ipStr != "" {
		rr.(*dns.A).A = net.ParseIP(ipStr).To4()
	} else {
		rr.(*dns.A).A = defaultIP
	}

	return rr
}

func dnsTXT(s string) *dns.TXT {
	t := new(dns.TXT)
	t.Txt = []string{"Client: " + s}
	t.Hdr = dns.RR_Header{
		Name:   *fqdn,
		Rrtype: dns.TypeTXT,
		Class:  dns.ClassINET,
		Ttl:    0,
	}

	return t
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
