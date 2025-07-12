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

const (
	IORING_OP_NOP uint8 = iota
	IORING_OP_READV
	IORING_OP_WRITEV
	IORING_OP_FSYNC
	IORING_OP_READ_FIXED
	IORING_OP_WRITE_FIXED
	IORING_OP_POLL_ADD
	IORING_OP_POLL_REMOVE
	IORING_OP_SYNC_FILE_RANGE
	IORING_OP_SENDMSG
	IORING_OP_RECVMSG
	IORING_OP_TIMEOUT
	IORING_OP_TIMEOUT_REMOVE
	IORING_OP_ACCEPT
	IORING_OP_ASYNC_CANCEL
	IORING_OP_LINK_TIMEOUT
	IORING_OP_CONNECT
	IORING_OP_FALLOCATE
	IORING_OP_OPENAT
	IORING_OP_CLOSE
	IORING_OP_FILES_UPDATE
	IORING_OP_STATX
	IORING_OP_READ
	IORING_OP_WRITE
	IORING_OP_FADVISE
	IORING_OP_MADVISE
	IORING_OP_SEND
	IORING_OP_RECV
	IORING_OP_OPENAT2
	IORING_OP_EPOLL_CTL
	IORING_OP_SPLICE
	IORING_OP_PROVIDE_BUFFERS
	IORING_OP_REMOVE_BUFFERS
	IORING_OP_TEE
	IORING_OP_SHUTDOWN
	IORING_OP_RENAMEAT
	IORING_OP_UNLINKAT
	IORING_OP_MKDIRAT
	IORING_OP_SYMLINKAT
	IORING_OP_LINKAT
	IORING_OP_MSG_RING
	IORING_OP_FSETXATTR
	IORING_OP_SETXATTR
	IORING_OP_FGETXATTR
	IORING_OP_GETXATTR
	IORING_OP_SOCKET
	IORING_OP_URING_CMD
	IORING_OP_SEND_ZC
	IORING_OP_SENDMSG_ZC
	IORING_OP_READ_MULTISHOT
	IORING_OP_WAITID
	IORING_OP_FUTEX_WAIT
	IORING_OP_FUTEX_WAKE
	IORING_OP_FUTEX_WAITV
	IORING_OP_FIXED_FD_INSTALL
	IORING_OP_FTRUNCATE
	IORING_OP_BIND
	IORING_OP_LISTEN
	IORING_OP_RECV_ZC
	IORING_OP_EPOLL_WAIT
	IORING_OP_READV_FIXED
	IORING_OP_WRITEV_FIXED
	IORING_OP_PIPE
	IORING_OP_LAST
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

type IOUringSQE struct {
	Opcode uint8 // type of opeartion for this sqe
	Flags  uint8 // IOSQE flags
	IOPrio uint16
	Fd     int32 // file descriptor to do io on

	Off  uint64 // first union: 1) off u64, 2) addr2 u64, 3) cmd_op u32, pad u32
	Addr uint64 // second union: 1) addr u64, 2) splice_off_in u64, 3) level u32 optname u32

	Len uint32

	Opflags uint32 // third union: a bunch of flags u32 or u16

	User_data uint64

	Buf_index   uint16 // fourth union: 1) buf_index u16, 2) buf_group u16
	Personality uint16
	Optlen      uint32 // fifth union: yeah im lazy i wont list them
	Optval      uint64 // sixth union: ditto
}

// fd is the returned fd of IOUringSetup
// toSubmit is how many sqe events we want to submit
func iOUringEnter(fd int, toSubmit uint32) error {
	_, _, e1 := unix.Syscall6(unix.SYS_IO_URING_ENTER, uintptr(toSubmit), uintptr(fd), uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if e1 == 0 {
		return nil
	}
	return e1
}

func IOUringSubmit(fd int, sqes []IOUringSQE) error {
	return iOUringEnter(fd, uint32(len(sqes)))
}
