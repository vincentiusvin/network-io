package lowlevel

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	IORING_FEAT_SINGLE_MMAP     = 1 << 0
	IORING_FEAT_NODROP          = 1 << 1
	IORING_FEAT_SUBMIT_STABLE   = 1 << 2
	IORING_FEAT_RW_CUR_POS      = 1 << 3
	IORING_FEAT_CUR_PERSONALITY = 1 << 4
	IORING_FEAT_FAST_POLL       = 1 << 5
	IORING_FEAT_POLL_32BITS     = 1 << 6
	IORING_FEAT_SQPOLL_NONFIXED = 1 << 7
	IORING_FEAT_EXT_ARG         = 1 << 8
	IORING_FEAT_NATIVE_WORKERS  = 1 << 9
	IORING_FEAT_RSRC_TAGS       = 1 << 10
	IORING_FEAT_CQE_SKIP        = 1 << 11
	IORING_FEAT_LINKED_FILE     = 1 << 12
	IORING_FEAT_REG_REG_RING    = 1 << 13
	IORING_FEAT_RECVSEND_BUNDLE = 1 << 14
	IORING_FEAT_MIN_TIMEOUT     = 1 << 15
	IORING_FEAT_RW_ATTR         = 1 << 16
	IORING_FEAT_NO_IOWAIT       = 1 << 17
)

const (
	IORING_OFF_SQ_RING    uint64 = 0
	IORING_OFF_CQ_RING    uint64 = 0x8000000
	IORING_OFF_SQES       uint64 = 0x10000000
	IORING_OFF_PBUF_RING  uint64 = 0x80000000
	IORING_OFF_PBUF_SHIFT        = 16
	IORING_OFF_MMAP_MASK  uint64 = 0xf8000000
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

const (
	SizeOfIOUringCQE = unsafe.Sizeof(IOUringCQE{})
	SizeOfUint32     = unsafe.Sizeof(uint32(0))
)

type IOUringCQE struct {
	User_data uint64
	Res       int32
	Flags     uint32
}

func IOUringSetup(entries uint32, params *IOUringParams) (fd int, err error) {
	paramsPtr := unsafe.Pointer(params)
	r0, _, e1 := unix.Syscall(unix.SYS_IO_URING_SETUP, uintptr(entries), uintptr(paramsPtr), 0)
	if e1 == 0 {
		return int(r0), nil
	}
	return int(r0), e1
}
