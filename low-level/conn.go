package lowlevel

import (
	"io"
	"syscall"
)

type SockFD int
type ConnFD int

func OpenSocket(port int) (SockFD, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return 0, err
	}

	sa := new(syscall.SockaddrInet4)
	sa.Port = port
	sa.Addr = [...]byte{0, 0, 0, 0}
	if err = syscall.Bind(fd, sa); err != nil {
		return 0, err
	}

	if err = syscall.Listen(fd, 10); err != nil {
		return 0, err
	}

	return SockFD(fd), nil
}

func (fd SockFD) Close() error {
	return syscall.Close(int(fd))
}

func (fd SockFD) AcceptConnection() (ConnFD, *syscall.SockaddrInet4, error) {
	// this will immediately return if we set syscall.SOCK_NONBLOCK above and we don't have anything queued.
	// fun!
	nfd, sa, err := syscall.Accept(int(fd))
	if err != nil {
		return 0, nil, err
	}
	sa4 := sa.(*syscall.SockaddrInet4)
	return ConnFD(nfd), sa4, nil
}

var _ io.ReadWriter = ConnFD(0)

func (cd ConnFD) Read(b []byte) (int, error) {
	return syscall.Read(int(cd), b)
}

func (cd ConnFD) Write(b []byte) (int, error) {
	return syscall.Write(int(cd), b)
}

func (cd ConnFD) Close() error {
	return syscall.Close(int(cd))
}
