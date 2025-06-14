package main

import (
	"golang.org/x/sys/unix"
)

func main() {
}

func epoll() {
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	defer unix.Close(epfd)

	ev := new(unix.EpollEvent)
	ev.Events = unix.EPOLLIN
	if err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, 0, ev); err != nil {
		panic(err)
	}

	num := 10
	evs := make([]unix.EpollEvent, num)
	for {
		_, err := unix.EpollWait(epfd, evs, 3000)
		if err != nil {
			break
		}
	}
}
