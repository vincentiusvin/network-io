package e2e_test

import (
	"net"
	"reflect"
	"testing"
)

const ADDR = ":8000"

func TestConnect(t *testing.T) {
	_, err := net.Dial("tcp", ADDR)
	if err != nil {
		t.Fatalf("cannot connect to %v", ADDR)
	}
}

func TestEcho(t *testing.T) {
	runTCPTest(t, []byte("12345"))
}

func TestConcurrency(t *testing.T) {
	t1 := runTCPTest(t, []byte("123456"))
	t2 := runTCPTest(t, []byte("123457"))

	<-t1
	<-t2
}

func runTCPTest(t *testing.T, inData []byte) (done chan struct{}) {
	done = make(chan struct{})

	go func() {
		defer close(done)
		c, err := net.Dial("tcp", ADDR)
		if err != nil {
			t.Errorf("cannot connect to %v", ADDR)
		}

		if _, err = c.Write(inData); err != nil {
			t.Error(err)
		}

		outBuf := make([]byte, 1024)
		n, err := c.Read(outBuf)
		if err != nil {
			t.Error(err)
		}
		outData := outBuf[:n]
		if !reflect.DeepEqual(outData, inData) {
			t.Errorf("wrong data: exp %v got %v", inData, outData)
		}
	}()

	return done
}
