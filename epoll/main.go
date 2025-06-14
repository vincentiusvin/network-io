package main

import (
	lowlevel "learn_io/low-level"
	"log"

	"golang.org/x/sys/unix"
)

func main() {
	epoll()
}

func epoll() {
	sock, err := lowlevel.OpenSocket(8000)
	if err != nil {
		panic(err)
	}
	sock.SetNonblock(true)
	defer sock.Close()

	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	defer unix.Close(epfd)

	ev := new(unix.EpollEvent)
	ev.Events = unix.EPOLLIN
	ev.Fd = int32(sock)
	if err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(sock), ev); err != nil {
		panic(err)
	}

	num := 1
	evs := make([]unix.EpollEvent, num)
	for {
		log.Printf("Starting to wait...")
		n, err := unix.EpollWait(epfd, evs, -1)
		if err != nil {
			break
		}
		log.Printf("Got %v events", n)
		for _, event := range evs[:n] {
			currFd := lowlevel.SockFD(event.Fd)
			if currFd != sock {
				continue
			}
			connFd, conn, err := currFd.AcceptConnection()
			if err != nil {
				panic(err)
			}
			log.Printf("Got connection %v:%v", conn.Addr, conn.Port)
			connFd.SetNonblock(true)
		}
	}
}
