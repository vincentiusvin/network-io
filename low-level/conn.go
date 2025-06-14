package lowlevel

import "syscall"

type SockFD int
type ConnFD int

func OpenSocket() (SockFD, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.SOCK_NONBLOCK, 0)
	if err != nil {
		return 0, err
	}

	sa := new(syscall.SockaddrInet4)
	sa.Port = 8000
	sa.Addr = [...]byte{0, 0, 0, 0}
	if err = syscall.Bind(fd, sa); err != nil {
		return 0, err
	}

	if err = syscall.Listen(fd, 10); err != nil {
		return 0, err
	}

	return SockFD(fd), nil
}

func (fd SockFD) AcceptConnection() (ConnFD, *syscall.SockaddrInet4, error) {
	// this will immediately return if we set syscall.SOCK|NONBLOCK above and we don't have anything queued.
	// fun!
	nfd, sa, err := syscall.Accept(int(fd))
	if err != nil {
		return 0, nil, err
	}
	sa4 := sa.(*syscall.SockaddrInet4)
	return ConnFD(nfd), sa4, nil
}
