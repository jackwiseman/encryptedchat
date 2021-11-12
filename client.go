package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":3540")
	if err != nil {
		fmt.Printf("Did not connect")
	} else {
		fmt.Printf("Connected, hooray!\n")
		//var test string
		//fmt.Scanf("%s", &test)
		conn.Write([]byte("This is a mesage from the client"))
	}

	b := bufio.NewReader(conn)

	for {
		line, err := b.ReadString('\n')
		if err != nil { // EOF, or worse
			break
		}
		//client.Write(line)
		fmt.Printf(line)
	}

	conn.Close()
}
