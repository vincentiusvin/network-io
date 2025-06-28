package lowlevel

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// see man io_uring_setup (2)

type IOUringParams struct {
	Sq_entries     uint32
	Cq_entries     uint32
	Flags          uint32
	Sq_thread_cpu  uint32
	Sq_thread_idle uint32
	Features       uint32
	Wq_fd          uint32
	Resv           [3]uint32
	Sq_off         IOSQringOffsets
	Cq_off         IOCQringOffsets
}

type IOSQringOffsets struct {
	Head         uint32
	Tail         uint32
	Ring_mask    uint32
	Ring_entries uint32
	Flags        uint32
	Dropped      uint32
	Array        uint32
	Resv         [3]uint32
}

type IOCQringOffsets struct {
	Head         uint32
	Tail         uint32
	Ring_mask    uint32
	Ring_entries uint32
	Overflow     uint32
	Cqes         uint32
	Flags        uint32
	Resv         [3]uint32
}

func IOUringEnter(entries uint32, params *IOUringParams) (fd int, err error) {
	paramsPtr := unsafe.Pointer(params)
	r0, _, e1 := unix.Syscall(unix.SYS_IO_URING_SETUP, uintptr(entries), uintptr(paramsPtr), 0)
	if e1 == 0 {
		return int(r0), nil
	}
	return int(r0), e1
}
