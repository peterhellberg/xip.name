package main

import (
	"net"
	"testing"
)

var clientStringTests = []struct {
	addr net.Addr
	out  string
}{
	{&net.UDPAddr{net.ParseIP("127.0.0.1"), 1234, ""}, "127.0.0.1:1234 (udp)"},
	{&net.TCPAddr{net.ParseIP("127.0.0.2"), 5678, ""}, "127.0.0.2:5678 (tcp)"},
}

func TestClientString(t *testing.T) {
	for _, tt := range clientStringTests {
		if got := clientString(tt.addr); got != tt.out {
			t.Errorf("clientSide(%#v) = %v, want %v", tt.addr, got, tt.out)
		}
	}
}
