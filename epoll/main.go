package main

import (
	"fmt"
	"io"
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

	if n == 0 {
		return nil
	}

	log.Printf("epoll wait returned %v events", n)

	for _, event := range evs[:n] {
		if event.Fd == int32(s.sockFD) {
			asSockFD := lowlevel.SockFD(event.Fd)
			if err := s.handleNewConnection(asSockFD); err != nil {
				return fmt.Errorf("cannot handle new connection: %v", err)
			}
		} else {
			isIn := event.Events&unix.EPOLLIN == unix.EPOLLIN
			isOut := event.Events&unix.EPOLLOUT == unix.EPOLLOUT
			isHup := event.Events&unix.EPOLLHUP == unix.EPOLLHUP
			msg := fmt.Sprintf("processing fd %v (old) events:", event.Fd)
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
					return fmt.Errorf("cannot handle existing connection: %v", err)
				}
			} else {
				log.Printf("closing fd %v", asConnFD)
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

// return non nil error if there is something that we cannot handle
func (s *server) handleExistingConnectionIn(connFD lowlevel.ConnFD) error {
	for {
		b, err := readConnection(connFD)
		if err != nil {
			if err == io.EOF {
				connFD.Close()
				return nil
			}
			if err == unix.EAGAIN {
				return nil
			}

			return fmt.Errorf("read error: %v", err)
		}

		if err := writeConnection(b, connFD); err != nil {
			return fmt.Errorf("write error: %v", err)
		}
	}
}

func readConnection(connFD lowlevel.ConnFD) ([]byte, error) {
	b := make([]byte, 1024)

	log.Printf("reading from fd %v", connFD)
	n, err := connFD.Read(b)

	if err != nil {
		log.Printf("reading fd %v error %v", connFD, err)
		return nil, err
	}
	if n == 0 {
		log.Printf("reading fd %v returns 0 bytes", connFD)
		return nil, io.EOF
	}
	log.Printf("read %v bytes from fd %v", n, connFD)
	return b[:n], nil
}

func writeConnection(b []byte, connFD lowlevel.ConnFD) error {
	log.Printf("writing to fd %v", connFD)
	n, err := connFD.Write(b)

	if err != nil {
		log.Printf("writing fd %v error %v", connFD, err)
		return err
	}
	if n == 0 {
		log.Printf("writing fd %v returns 0 bytes", connFD)
		return nil
	}
	log.Printf("written %v bytes to fd %v", n, connFD)
	return nil
}

func (s *server) handleNewConnection(sockFd lowlevel.SockFD) error {
	connFD, conn, err := sockFd.AcceptConnection()
	if err != nil {
		return err
	}
	connFD.SetNonblock(true)
	log.Printf("got new connection from port %v: allocating to fd %v", conn.Port, connFD)
	s.registerFDToEpoll(int(connFD), unix.EPOLLIN|unix.EPOLLOUT|unix.EPOLLHUP)
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
