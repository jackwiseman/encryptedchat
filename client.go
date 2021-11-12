package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":3540")
	if err != nil {
		fmt.Printf("Did not connect")
	} else {
		fmt.Printf("Connected, hooray!\n")
		conn.Write([]byte("This is a message from the client\n"))
	}
	conn.Close()
}
