package main

import lowlevel "learn_io/low-level"

func main() {
	fd, err := lowlevel.OpenSocket()
	if err != nil {
		panic(err)
	}
	print(fd)
}
