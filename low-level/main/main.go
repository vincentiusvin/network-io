package main

import (
	lowlevel "learn_io/low-level"
	"log"
)

func main() {
	fd, err := lowlevel.OpenSocket(8000)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	nfd, conn, err := fd.AcceptConnection()
	if err != nil {
		panic(err)
	}
	defer nfd.Close()

	log.Printf("Got connection from: %v:%v", conn.Addr, conn.Port)
	for {
		b := make([]byte, 1024)
		n, err := nfd.Read(b)
		if err != nil {
			break
		}
		log.Print(string(b[:n]))
	}
}
