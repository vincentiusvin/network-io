package main

import (
	lowlevel "learn_io/low-level"
	"log"

	"golang.org/x/sys/unix"
)

func main() {
	s := NewServer(8000)
	s.Listen()
}

type server struct {
	port      int
	connQueue int

	sockFD   lowlevel.SockFD
	sockEpFD int
}

func NewServer(port int) *server {
	return &server{
		port:      8000,
		connQueue: 1,
	}
}

func (s *server) Listen() error {
	sock, err := lowlevel.OpenSocket(8000)
	if err != nil {
		return err
	}
	sock.SetNonblock(true)
	s.sockFD = sock
	s.pollForConnections()
	return nil
}

func (s *server) pollForConnections() error {
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return err
	}
	s.sockEpFD = epfd

	ev := new(unix.EpollEvent)
	ev.Events = unix.EPOLLIN
	ev.Fd = int32(s.sockFD)
	if err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(s.sockFD), ev); err != nil {
		return err
	}

	num := 1
	evs := make([]unix.EpollEvent, num)
	for {
		log.Printf("Starting to wait...")
		n, err := unix.EpollWait(s.sockEpFD, evs, -1)
		if err != nil {
			break
		}

		log.Printf("Got %v events", n)

		for _, event := range evs[:n] {
			currFd := lowlevel.SockFD(event.Fd)
			if currFd != s.sockFD {
				continue
			}
			s.handleNewConnection(currFd)
		}
	}
	return nil
}

func (s *server) handleNewConnection(sockFd lowlevel.SockFD) {
	connFd, conn, err := sockFd.AcceptConnection()
	if err != nil {
		panic(err)
	}

	log.Printf("Got connection %v:%v", conn.Addr, conn.Port)

	connFd.SetNonblock(true)
}
