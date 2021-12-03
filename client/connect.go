package main

import (
	//"crypto/rand"
	"crypto/rsa"
	//"crypto/sha256"
	//"encoding/base64"
	"encoding/gob"
	//"strconv"
	//"time"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	Msg string
	Publickey rsa.PublicKey
}

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("# Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	gob.Register(new(Message))
	dec := gob.NewDecoder(c)
	enc := gob.NewEncoder(c)

	ch := make(chan Message)
	fmt.Println("# You have joined encryptedchat. Type /help for more info, /quit to exit.")
	go printer(ch)
	go listener(ch, dec)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text == ""{
			continue
		}
		msg := new(Message)
		msg.Msg = text

		err := enc.Encode(msg)
		if err != nil {
			panic(err)
		}

		if strings.TrimSpace(string(text)) == "/quit" {
			fmt.Println("# Disconnected")
			return
		}
	}
}
func printer(ch chan Message) {
	for {
		msg := <- ch
		fmt.Printf(msg.Msg)
	}
}

func listener(ch chan Message, dec *gob.Decoder){
	for {
		msg := new(Message)
		err := dec.Decode(msg)
		if err != nil {
			panic(err)
		}
		ch <- *msg
	}
}
