package main

import (
	lowlevel "learn_io/low-level"

	"golang.org/x/sys/unix"
)

func main() {
	s := NewServer(8000)
	if err := s.Listen(); err != nil {
		panic(err)
	}
	for {
		if err := s.Loop(1, -1); err != nil {
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

func (s *server) Loop(queueSize int, timeout int) error {
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
			s.handleExistingConnection(asConnFD)
		}
	}

	return nil
}

func (s *server) handleExistingConnection(connFd lowlevel.ConnFD) {

}

func (s *server) handleNewConnection(sockFd lowlevel.SockFD) {
	connFd, _, err := sockFd.AcceptConnection()
	if err != nil {
		panic(err)
	}
	connFd.SetNonblock(true)

	s.registerFDToEpoll(int(sockFd), unix.EPOLL_CTL_ADD)
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
