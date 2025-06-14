package main

import (
	"log"
	"net"
)

func main() {
	port := ":8000"

	list, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	log.Printf("Listening at %v", port)

	for {
		conn, err := list.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	log.Printf("Got conn from %v", c.RemoteAddr())
	defer c.Close()
	for {
		b := make([]byte, 1024)

		n, err := c.Read(b)
		if err != nil {
			log.Println(err)
			break
		}

		_, err = c.Write(b[:n])
		if err != nil {
			log.Println(err)
			break
		}
	}
}
