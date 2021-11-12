// $ ./echo
// $ nc localhost 3540

package main

import (
	"fmt"
	"net"
	"strconv"
	"io"
)

const PORT = 3540

func reciever(ch chan string) {
	for {
		fmt.Println(<-ch)
	}
}

func main() {

	server, err := net.Listen("tcp", ":" + strconv.Itoa(PORT))
	if server == nil {
		panic("couldn't start listening: " + err.Error())
	}
	
	ch := make(chan string)

	conns := clientConns(server)
	for {
		go reciever(ch)
		go handleConn(<-conns, ch)
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

func handleConn(client net.Conn, ch chan string) {

	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 256)

	for {
		n, err := client.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error")
			}
			break
		}
		buf = append(buf, tmp[:n]...)
		ch <- string(buf)
	}

}

func response(client net.Conn, line string) {
	out := "echo: " + line
	client.Write([]byte(out))
}
