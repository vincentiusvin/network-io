package main

import (
	lowlevel "learn_io/low-level"
	"log"
	"os"
	"os/signal"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fd.Close()
		nfd.Close()
		os.Exit(1)
	}()

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
