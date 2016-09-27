package main

import (
	"net"
	"testing"

	"github.com/miekg/dns"
)

func TestHandleDNS(t *testing.T) {
	for _, tt := range []struct {
		w dns.ResponseWriter
		r *dns.Msg
		n string
	}{
		{&fakeResponseWriter{}, &dns.Msg{Question: []dns.Question{
			dns.Question{"c7.se.", dns.TypeMX, dns.ClassINET},
		}}, "c7.se."},
		{&fakeResponseWriter{}, &dns.Msg{Question: []dns.Question{
			dns.Question{"xip.name.", dns.TypeAAAA, dns.ClassINET},
		}}, "xip.name."},
	} {
		handleDNS(tt.w, tt.r)

		f, ok := tt.w.(*fakeResponseWriter)
		if !ok {
			t.Fatalf("tt.w is not a *fakeResponseWriter")
		}

		if got, want := f.Msg.Answer[0].Header().Name, tt.n; got != want {
			t.Fatalf("f.Msg.Answer[0].Header().Name = %q, want %q", got, want)
		}
	}
}

func TestNewServer(t *testing.T) {
	for _, tt := range []struct {
		net string
	}{
		{"abc"},
		{"xyz"},
	} {
		s := newServer(tt.net)

		if got, want := s.Net, tt.net; got != want {
			t.Fatalf("s.Net = %q, want %q", got, want)
		}
	}
}

func TestDnsRR(t *testing.T) {
	for _, tt := range []struct {
		name string
		want string
	}{
		{"abc", "abc\t300\tIN\tA\t"},
		{"xyz", "xyz\t300\tIN\tA\t"},
	} {
		rr := dnsRR(tt.name)

		if got, want := rr.String(), tt.want; got != want {
			t.Fatalf("rr.String() = %q, want %q", got, want)
		}
	}
}

func TestDnsTXT(t *testing.T) {
	for _, tt := range []struct {
		s   string
		txt string
	}{
		{"abc", "Client: abc"},
		{"xyz", "Client: xyz"},
	} {
		dt := dnsTXT(tt.s)

		if got, want := len(dt.Txt), 1; got != want {
			t.Fatalf("len(dt.Txt) = %d, want %d", got, want)
		}

		if got, want := dt.Txt[0], tt.txt; got != want {
			t.Fatalf("dt.Txt[0] = %q, want %q", got, want)
		}

		if got, want := dt.Hdr.Rrtype, dns.TypeTXT; got != want {
			t.Fatalf("dt.Hdr.Rrtype = %d, want %d", got, want)
		}
	}
}

func TestClientString(t *testing.T) {
	for _, tt := range []struct {
		addr net.Addr
		out  string
	}{
		{unknownAddr{}, "unknown"},
		{&net.UDPAddr{net.ParseIP("127.0.0.1"), 1234, ""}, "127.0.0.1:1234 (udp)"},
		{&net.TCPAddr{net.ParseIP("127.0.0.2"), 5678, ""}, "127.0.0.2:5678 (tcp)"},
	} {
		if got := clientString(tt.addr); got != tt.out {
			t.Errorf("clientSide(%#v) = %v, want %v", tt.addr, got, tt.out)
		}
	}
}

type fakeResponseWriter struct {
	Msg *dns.Msg
}

func (f *fakeResponseWriter) LocalAddr() net.Addr {
	panic("not implemented LocalAddr")
}

func (f *fakeResponseWriter) RemoteAddr() net.Addr {
	return &net.TCPAddr{net.ParseIP("127.0.0.1"), 5678, ""}
}

func (f *fakeResponseWriter) WriteMsg(msg *dns.Msg) error {
	f.Msg = msg
	return nil
}

func (f *fakeResponseWriter) Write([]byte) (int, error) {
	panic("Write not implemented")
}

func (f *fakeResponseWriter) Close() error {
	panic("Close not implemented")
}

func (f *fakeResponseWriter) TsigStatus() error {
	panic("TsigStatus not implemented")
}

func (f *fakeResponseWriter) TsigTimersOnly(bool) {
	panic("TsigTimersOnly not implemented")
}

func (f *fakeResponseWriter) Hijack() {
	panic("Hijack not implemented")
}

type unknownAddr struct{}

func (unknownAddr) Network() string {
	return ""
}

func (unknownAddr) String() string {
	return ""
}
