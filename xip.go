// Copyright 2014 Peter Hellberg.
// Released under the terms of the MIT license.

// xip is a small name server which sends back any IP address found in the provided hostname.
//
// When queried for type A, it sends back the parsed IPv4 address.
// In the additional section the port number and transport are shown.
//
// Basic use pattern:
//
// 		dig @xip.name www.xip.name A
//
// 		; <<>> DiG 9.8.3-P1 <<>> @xip.name www.xip.name A
// 		; (1 server found)
// 		;; global options: +cmd
// 		;; Got answer:
// 		;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 47078
// 		;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
// 		;; WARNING: recursion requested but not available
//
// 		;; QUESTION SECTION:
// 		;www.xip.name.			IN	A
//
// 		;; ANSWER SECTION:
// 		www.xip.name.		0	IN	A	188.166.43.179
//
// 		;; ADDITIONAL SECTION:
// 		xip.name.		0	IN	TXT	"IP: 188.126.74.76:58956 (udp)"
//
// 		;; Query time: 31 msec
// 		;; SERVER: 188.166.43.179#53(188.166.43.179)
// 		;; WHEN: Wed Dec 31 02:13:50 2014
// 		;; MSG SIZE  rcvd: 108
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
	verbose  = flag.Bool("v", false, "Verbose")
	compress = flag.Bool("c", false, "compress replies")
	fqdn     = flag.String("fqdn", "xip.name.", "FQDN to handle")
	port     = flag.String("p", "53", "The port to bind on")
	ip       = flag.String("ip", "188.166.43.179", "The IP of xip.name")

	ipPattern = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
)

func main() {
	flag.Parse()

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
	m.Compress = *compress

	var (
		rr  dns.RR
		str string
	)

	if ip, ok := w.RemoteAddr().(*net.UDPAddr); ok {
		str = "IP: " + ip.String() + " (udp)"
	}

	if ip, ok := w.RemoteAddr().(*net.TCPAddr); ok {
		str = "IP: " + ip.String() + " (tcp)"
	}

	rr = new(dns.A)
	rr.(*dns.A).Hdr = dns.RR_Header{
		Name:   r.Question[0].Name,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    0,
	}

	if r.Question[0].Name == "xip.name." || r.Question[0].Name == "www.xip.name." {
		rr.(*dns.A).A = net.ParseIP("188.166.43.179").To4()
	} else {
		ipStr := ipPattern.FindString(r.Question[0].Name)

		rr.(*dns.A).A = net.ParseIP(ipStr).To4()
	}

	t := new(dns.TXT)
	t.Hdr = dns.RR_Header{
		Name:   *fqdn,
		Rrtype: dns.TypeTXT,
		Class:  dns.ClassINET,
		Ttl:    0,
	}
	t.Txt = []string{str}

	switch r.Question[0].Qtype {
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

		soa, _ := dns.NewRR(`xip.name. 0 IN SOA 2009032802 21600 7200 604800 3600`)

		c <- &dns.Envelope{RR: []dns.RR{soa, t, rr, soa}}
		w.Hijack()
		// w.Close() // Client closes connection
		return

	}

	if *verbose {
		fmt.Printf("%v\n", m.String())
	}

	w.WriteMsg(m)
}

func serve(net string) {
	server := &dns.Server{
		Addr:       ":" + *port,
		Net:        net,
		TsigSecret: nil,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())
	}
}
