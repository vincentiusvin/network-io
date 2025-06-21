package main

import (
	lowlevel "learn_io/low-level"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
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
		nfd, err := s.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(nfd)
	}
}

func handleConnection(connFD lowlevel.ConnFD) {
	defer connFD.Close()

	for {
		b := make([]byte, 1024)
		waiting.Store(connFD, "reading")
		n, err := connFD.Read(b)
		waiting.Delete(connFD)
		if err != nil || n == 0 {
			return
		}
		waiting.Store(connFD, "writing")
		n, err = connFD.Write(b[:n])
		waiting.Delete(connFD)
		if err != nil || n == 0 {
			return
		}
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

func (s *server) Accept() (lowlevel.ConnFD, error) {
	nfd, _, err := s.sockFD.AcceptConnection()
	if err != nil {
		return lowlevel.ConnFD(0), err
	}
	return nfd, nil
}

func (s *server) Close() {
	s.sockFD.Close()
}
