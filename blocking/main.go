package main

import (
	lowlevel "learn_io/low-level"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

var waiting sync.Map

func main() {
	s := NewServer(8000)
	if err := s.Listen(); err != nil {
		panic(err)
	}
	defer s.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		s.Close()
		os.Exit(1)
	}()

	go func() {
		for {
			log.Print("routines:", runtime.NumGoroutine())
			waiting.Range(
				func(key, value any) bool {
					log.Printf("fd %v status %v", key, value)
					return true
				},
			)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		nfd, info, err := s.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(nfd, info)
	}
}

func handleConnection(connFD lowlevel.ConnFD, info *unix.SockaddrInet4) {
	log.Printf("got new connection from port %v: allocating to fd %v", info.Port, connFD)

	defer func() {
		log.Printf("closing fd %v", connFD)
		connFD.Close()
	}()

	for {
		b := make([]byte, 1024)

		waiting.Store(connFD, "reading")
		log.Printf("reading from fd %v", connFD)
		n, err := connFD.Read(b)
		waiting.Delete(connFD)

		if err != nil {
			log.Printf("reading fd %v error %v", connFD, err)
			return
		}
		if n == 0 {
			log.Printf("reading fd %v returns 0 bytes", connFD)
			return
		}
		log.Printf("read %v bytes from fd %v", n, connFD)

		waiting.Store(connFD, "writing")
		log.Printf("writing to fd %v", connFD)
		n, err = connFD.Write(b[:n])
		waiting.Delete(connFD)

		if err != nil {
			log.Printf("writing fd %v error %v", connFD, err)
			return
		}
		if n == 0 {
			log.Printf("writing fd %v returns 0 bytes", connFD)
			return
		}
		log.Printf("written %v bytes to fd %v", n, connFD)
	}
}

type server struct {
	port   int
	sockFD lowlevel.SockFD
}

func NewServer(port int) *server {
	return &server{
		port: port,
	}
}

func (s *server) Listen() error {
	fd, err := lowlevel.OpenSocket(s.port)
	if err != nil {
		return err
	}
	s.sockFD = fd
	return nil
}

func (s *server) Accept() (lowlevel.ConnFD, *unix.SockaddrInet4, error) {
	return s.sockFD.AcceptConnection()
}

func (s *server) Close() error {
	return s.sockFD.Close()
}
