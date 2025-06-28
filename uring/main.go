package main

import (
	lowlevel "learn_io/low-level"
	"log"

	"golang.org/x/sys/unix"
)

func main() {
	params := new(lowlevel.IOUringParams)
	fd, err := lowlevel.IOUringSetup(10, params)
	if err != nil {
		panic(err)
	}

	sq, err := mmapSQ(fd, params)
	if err != nil {
		panic(err)
	}
	cq, err := mmapCQ(fd, params)
	if err != nil {
		panic(err)
	}
	log.Println(sq, cq)
}

func mmapSQ(ringFD int, params *lowlevel.IOUringParams) ([]byte, error) {
	sqRingLen := params.Sq_off.Array + params.Sq_entries*uint32(lowlevel.SizeOfUint32)

	// Go's mmap(2) is different from C so be careful!
	return unix.Mmap(
		ringFD,
		int64(lowlevel.IORING_OFF_SQ_RING),
		int(sqRingLen),
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_SHARED|unix.MAP_POPULATE,
	)
}

func mmapCQ(ringFD int, params *lowlevel.IOUringParams) ([]byte, error) {
	cqRingLen := params.Cq_off.Cqes + params.Cq_entries*uint32(lowlevel.SizeOfIOUringCQE)

	// Go's mmap(2) is different from C so be careful!
	return unix.Mmap(
		ringFD,
		int64(lowlevel.IORING_OFF_CQ_RING),
		int(cqRingLen),
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_SHARED|unix.MAP_POPULATE,
	)
}
