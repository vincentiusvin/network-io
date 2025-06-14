package main

import (
	lowlevel "learn_io/low-level"

	"golang.org/x/sys/unix"
)

func main() {
	s := NewServer(8000)
	s.Listen()
}

type server struct {
	port int

	sockFD   lowlevel.SockFD
	sockEpFD int
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

	epfd, err := setupEpoll(int(sock), unix.EPOLLIN)
	if err != nil {
		return err
	}
	s.sockEpFD = epfd
	return nil
}

func (s *server) Accept() error {
	connQueue := 10
	evs := make([]unix.EpollEvent, connQueue)

	n, err := unix.EpollWait(s.sockEpFD, evs, 10)
	if err != nil {
		return err
	}

	for _, event := range evs[:n] {
		currFd := lowlevel.SockFD(event.Fd)
		if currFd != s.sockFD {
			continue
		}
		s.handleNewConnection(currFd)
	}

	return nil
}

func (s *server) handleNewConnection(sockFd lowlevel.SockFD) {
	connFd, _, err := sockFd.AcceptConnection()
	if err != nil {
		panic(err)
	}
	connFd.SetNonblock(true)

	epfd, err := setupEpoll(int(connFd), unix.EPOLLIN)
}

func setupEpoll(fd int, events uint32) (int, error) {
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return 0, err
	}

	ev := new(unix.EpollEvent)
	ev.Events = events
	ev.Fd = int32(fd)
	if err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, fd, ev); err != nil {
		return 0, err
	}
	return epfd, nil
}
