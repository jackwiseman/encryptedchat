package main

import (
	"net"
	"fmt"
)

func main() {
	conn, err := net.Dial("tcp", ":3540")
	if err != nil {
		fmt.Printf("Did not connect")	
	} else {
		fmt.Printf("Connected, hooray!")
	}
	conn.Close()
}
