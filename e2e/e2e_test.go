package e2e_test

import (
	"bytes"
	"crypto/rand"
	"io"
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

func TestLotsOfData(t *testing.T) {
	b := make([]byte, 10240000)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		t.Fatal(err)
	}

	t1 := runTCPTest(t, b)

	<-t1
}

func runTCPTest(t *testing.T, inData []byte) (done chan struct{}) {
	done = make(chan struct{})

	go func() {
		defer close(done)
		c, err := net.Dial("tcp", ADDR)
		if err != nil {
			t.Errorf("cannot connect to %v", ADDR)
		}

		inReader := bytes.NewReader(inData)
		if _, err = io.Copy(c, inReader); err != nil {
			t.Error(err)
		}

		outData := make([]byte, len(inData))
		n, err := io.ReadFull(c, outData)
		if err != nil {
			t.Error(err)
		}
		outData = outData[:n]

		if !reflect.DeepEqual(outData, inData) {
			t.Errorf("wrong data: exp %v(len %v) got %v(len %v)", inData[:5], len(inData), outData[:5], len(outData))
		}
	}()

	return done
}
