package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
)

// TODO: Next two functions

/*func displayMessage(ch chan string) {
	for {
		fmt.Println(<-ch)
	}
}

func listener(server net.Conn, ch chan string) {
	b := bufio.NewReader(server)
	for { // should only break when client disconnects
		line, err := b.ReadString('\n')
		if err != nil {
			break
		}
		ch <- line
	}
}*/

func sendToServer(server net.Conn, msg []byte) {
	server.Write([]byte(msg))
	fmt.Printf("To server -> " + string(msg) + "\n")
}

func main() {
	// Handle connection
	conn, err := net.Dial("tcp", ":3540")
	if err != nil {
		fmt.Printf("Did not connect")
	} else {
		fmt.Printf("Connected to server\n")
	}

	// This next section is mostly unneccessary for now, deals with
	// being able to have user input while messages are coming in

	var input []byte
	reader := bufio.NewReader(os.Stdin)
	// Disable input buffering
	exec.Command("stty", "-F", "/dev/tty/", "cbreak", "min", "1").Run()
	
	// Append each char to input byte

	// Keeps connection open indefinitely so user can send msgs
	for {
		b, err := reader.ReadByte()
		if err != nil {
			panic(err)
		} else {
			if b == 10 {// newline
				sendToServer(conn, input)
				input = nil
			} else {
				input = append(input, b)
			}
		}
	}

	conn.Close()
	fmt.Printf("Disconnected from server ")
}
