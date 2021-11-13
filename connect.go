package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch := make(chan string)
	fmt.Println("You have joined encryptedchat. Type /help for more info.")
	go printer(ch)
	go listener(ch, c)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")



		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}


	}
}
func printer(ch chan string) {
	for {
		fmt.Printf(<-ch)
	}
}

func listener(ch chan string, c net.Conn){
	for {
		message, _ := bufio.NewReader(c).ReadString('\n')
		ch <- message
	}
}