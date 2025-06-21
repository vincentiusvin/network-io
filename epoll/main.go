package main

import (
	"fmt"
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
		if err := s.Process(100, 100); err != nil {
			panic(err)
		}
	}
}

type server struct {
	port int

	sockFD lowlevel.SockFD
	epfd   int
}

func NewServer(port int) *server {
	return &server{
		port: 8000,
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
		if err == unix.EINTR {
			return nil
		}
		return err
	}

	log.Printf("Epoll wait returned %v events", n)

	if n == 0 {
		return nil
	}

	for _, event := range evs[:n] {
		if event.Fd == int32(s.sockFD) {
			log.Printf("Processing fd %v (new)", event.Fd)
			asSockFD := lowlevel.SockFD(event.Fd)
			err := s.handleNewConnection(asSockFD)
			if err != nil {
				return err
			}
		} else {
			isIn := event.Events&unix.EPOLLIN == unix.EPOLLIN
			isOut := event.Events&unix.EPOLLOUT == unix.EPOLLOUT
			isHup := event.Events&unix.EPOLLHUP == unix.EPOLLHUP
			msg := fmt.Sprintf("Processing fd %v (old) events:", event.Fd)
			if isIn {
				msg += " in"
			}
			if isOut {
				msg += " out"
			}
			if isHup {
				msg += " hup"
			}
			log.Print(msg)

			asConnFD := lowlevel.ConnFD(event.Fd)
			if isIn {
				err := s.handleExistingConnectionIn(asConnFD)
				if err != nil {
					return err
				}
			} else {
				asConnFD.Close()
			}
		}
	}

	return nil
}

func (s *server) Close() error {
	log.Println("Server closed")
	return s.sockFD.Close()
}

func (s *server) handleExistingConnectionIn(connFd lowlevel.ConnFD) error {
	for {
		b := make([]byte, 1024)
		read, err := connFd.Read(b)
		log.Printf("fd %v read %v bytes", connFd, read)
		if read == 0 {
			connFd.Close()
			return nil
		}

		if err != nil {
			if err == unix.EAGAIN {
				return nil
			}
			return err
		}

		written, err := connFd.Write(b[:read])
		log.Printf("fd %v write %v bytes", connFd, read)
		if written == 0 {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Sent %v bytes", written)
	}
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
