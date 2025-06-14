package main

import (
	"syscall"
)

func main() {
}

func epoll() {
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(epfd)

	ev := new(syscall.EpollEvent)
	ev.Events = syscall.EPOLLIN
	if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, 0, ev); err != nil {
		panic(err)
	}

	num := 10
	evs := make([]syscall.EpollEvent, num)
	for {
		_, err := syscall.EpollWait(epfd, evs, 3000)
		if err != nil {
			break
		}
	}
}
