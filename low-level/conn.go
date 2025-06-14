package lowlevel

import (
	"io"

	"golang.org/x/sys/unix"
)

type SockFD int
type ConnFD int

func OpenSocket(port int) (SockFD, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return 0, err
	}

	sa := new(unix.SockaddrInet4)
	sa.Port = port
	sa.Addr = [...]byte{0, 0, 0, 0}
	if err = unix.Bind(fd, sa); err != nil {
		return 0, err
	}

	if err = unix.Listen(fd, 10); err != nil {
		return 0, err
	}

	return SockFD(fd), nil
}

func (fd SockFD) Close() error {
	return unix.Close(int(fd))
}

func (fd SockFD) SetNonblock(nb bool) error {
	return unix.SetNonblock(int(fd), nb)
}

func (fd SockFD) AcceptConnection() (ConnFD, *unix.SockaddrInet4, error) {
	// this will immediately return if we set unix.SOCK_NONBLOCK above and we don't have anything queued.
	// fun!
	nfd, sa, err := unix.Accept(int(fd))
	if err != nil {
		return 0, nil, err
	}
	sa4 := sa.(*unix.SockaddrInet4)
	return ConnFD(nfd), sa4, nil
}

var _ io.ReadWriter = ConnFD(0)

func (cd ConnFD) Read(b []byte) (int, error) {
	return unix.Read(int(cd), b)
}

func (cd ConnFD) Write(b []byte) (int, error) {
	return unix.Write(int(cd), b)
}

func (cd ConnFD) Close() error {
	return unix.Close(int(cd))
}

func (cd ConnFD) SetNonblock(nb bool) error {
	return unix.SetNonblock(int(cd), nb)
}
