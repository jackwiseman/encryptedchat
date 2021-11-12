// $ ./echo
// $ nc localhost 3540

package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

const PORT = 3540

func main() {
	server, err := net.Listen("tcp", ":" + strconv.Itoa(PORT))
	if server == nil {
		panic("couldn't start listening: " + err.Error())
	}
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept: " + err.Error())
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	for {
		b := bufio.NewReader(client)
		line, err := b.ReadString('\n')
		if err != nil { // EOF, or worse
			break
		}
		//client.Write(line)
		fmt.Printf(line)
		response(client, line)
	}
}

func response(client net.Conn, line string) {
	out := "echo: " + line
	client.Write([]byte(out))
}