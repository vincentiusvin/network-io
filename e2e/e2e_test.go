package e2e_test

import (
	"net"
	"testing"
)

const ADDR = ":8000"

func TestConnect(t *testing.T) {
	_, err := net.Dial("tcp", ADDR)
	if err != nil {
		t.Fatalf("cannot connect to %v", ADDR)
	}
}
