package lowlevel

import "syscall"

func OpenSocket() (fd int, err error) {
	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.SOCK_NONBLOCK, 0)
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

	return fd, nil
}
