package main

import (
	lowlevel "learn_io/low-level"
	"log"
)

func main() {
	params := new(lowlevel.IOUringParams)
	fd, err := lowlevel.IOUringEnter(10, params)
	if err != nil {
		panic(err)
	}
	log.Print(fd)
}
