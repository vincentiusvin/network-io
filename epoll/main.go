package main

import (
	lowlevel "learn_io/low-level"
	"log"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

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

	for {
		if err := s.Process(1, -1); err != nil {
			panic(err)
		}
	}
}

type server struct {
	port int

	sockFD  lowlevel.SockFD
	epfd    int
	toWrite map[lowlevel.ConnFD][]byte
}

func NewServer(port int) *server {
	return &server{
		port:    8000,
		toWrite: make(map[lowlevel.ConnFD][]byte),
	}
}

func (s *server) Listen() error {
	sock, err := lowlevel.OpenSocket(8000)
	if err != nil {
		return err
	}
	sock.SetNonblock(true)
	s.sockFD = sock

	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return err
	}
	s.epfd = epfd

	s.registerFDToEpoll(int(s.sockFD), unix.EPOLLIN)

	return nil
}

func (s *server) Process(queueSize int, timeout int) error {
	evs := make([]unix.EpollEvent, queueSize)

	n, err := unix.EpollWait(s.epfd, evs, timeout)
	if err != nil {
		return err
	}

	for _, event := range evs[:n] {
		if event.Fd == int32(s.sockFD) {
			asSockFD := lowlevel.SockFD(event.Fd)
			s.handleNewConnection(asSockFD)
		} else {
			asConnFD := lowlevel.ConnFD(event.Fd)
			if event.Events&unix.EPOLLIN == unix.EPOLLIN {
				s.handleExistingConnection(asConnFD)
			}
		}
	}

	return nil
}

func (s *server) Close() error {
	return s.sockFD.Close()
}

func (s *server) handleExistingConnection(connFd lowlevel.ConnFD) error {
	b := make([]byte, 1024)
	read, err := connFd.Read(b)
	if read == 0 || err != nil {
		return err
	}
	log.Printf("Read %v bytes", read)

	written, err := connFd.Write(b[:read])
	if read == 0 || err != nil {
		return err
	}
	log.Printf("Sent %v bytes", written)
	if written < read {
		s.toWrite[connFd] = append(s.toWrite[connFd], b[written:read]...)
	}
	return nil
}

func (s *server) handleNewConnection(sockFd lowlevel.SockFD) error {
	connFd, conn, err := sockFd.AcceptConnection()
	if err != nil {
		return err
	}
	connFd.SetNonblock(true)
	log.Printf("conn %v:%v", conn.Addr, conn.Port)
	s.registerFDToEpoll(int(connFd), unix.EPOLLIN|unix.EPOLLOUT|unix.EPOLLHUP)
	return nil
}

func (s *server) registerFDToEpoll(fd int, events uint32) error {
	ev := new(unix.EpollEvent)
	ev.Events = events
	ev.Fd = int32(fd)
	if err := unix.EpollCtl(s.epfd, unix.EPOLL_CTL_ADD, fd, ev); err != nil {
		return err
	}
	return nil
}
