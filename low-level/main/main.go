package main

import (
	"fmt"
	lowlevel "learn_io/low-level"
)

func main() {
	fd, err := lowlevel.OpenSocket()
	if err != nil {
		panic(err)
	}
	nfd, conn, err := lowlevel.AcceptConnection(fd)
	if err != nil {
		panic(err)
	}
	fmt.Println(nfd, conn)
}
