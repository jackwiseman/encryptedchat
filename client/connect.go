package main

import (
	"crypto/rsa"
	"encoding/gob"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
//	"math/big"
)

type Message struct {
	Msg string
	Sender string
	PublicKey rsa.PublicKey
}

var encryptionKey rsa.PublicKey
var myKey rsa.PrivateKey
var serverKey rsa.PublicKey

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

	// Create and send key to server
	loadPrivateKey(&myKey)

	msg := new(Message)
	msg.PublicKey = myKey.PublicKey

	err = enc.Encode(msg)
	if err != nil {
		panic(err)
	}

	ch := make(chan Message)
	quit := make(chan bool, 2)

	fmt.Println("# You have joined encryptedchat. Type /help for more info, /quit to exit.")
	go printer(ch, quit, enc)
	go listener(ch, quit, dec)

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		select {
		case <- quit:
			return
		default:
			if input == ""{
				continue
			}
			msg := new(Message)
			if encryptionKey.E != 0 && string(input)[0:1] != "/" {
				msg.Msg = encrypt(input, encryptionKey)

			} else {
				msg.Msg = encrypt(input, serverKey)
				msg.Sender = "cmd"
			}

			err := enc.Encode(msg)
			if err != nil {
				panic(err)
			}

			if strings.TrimSpace(string(input)) == "/quit" {
				fmt.Println("# Disconnected")
				return
			}
		}
	}
}
func printer(ch chan Message, quit chan bool, enc *gob.Encoder) {
	for {
		msg := <-ch
		if msg.Sender == "auth" { // select statement
			// if we get the auth token, decrypt and send it back
			if msg.PublicKey.E != myKey.PublicKey.E || msg.PublicKey.N.Cmp(myKey.PublicKey.N) != 0 {
				fmt.Println("# Unable to authenticate, press enter to disconnect")
				quit <- true
				return
			}

			// this is likely the same thing as above, but safe to still error handle
			token, issue := decrypt(msg.Msg, myKey)
			if issue != nil {
				if strings.Contains(issue.Error(), "Decoding error") {
					panic(issue)
				}
				if strings.Contains(issue.Error(), "Decryption error") {
					fmt.Println("# Unable to authenticate, press enter to disconnect")
					quit <- true
					return
				}
			}

			response := new(Message)
			response.Sender = "auth"
			response.Msg = token
			err := enc.Encode(response)
			if err != nil {
				panic(err)
			}
			continue
		} else if msg.Sender == "serverkey" {
			serverKey = msg.PublicKey
			continue
		} else if msg.Sender == "server" {
			decrypted, _ := decrypt(msg.Msg, myKey)
			fmt.Printf(decrypted)
			continue
		} else if msg.Msg == "" {
			encryptionKey = msg.PublicKey
			continue
		} else {
			if msg.PublicKey.E != 0 {
				encryptionKey = msg.PublicKey
				decrypted, _ := decrypt(msg.Msg, myKey)
				fmt.Printf(msg.Sender + ": " + decrypted)
				fmt.Printf(decrypted)
			}
		}
	}
}

func listener(ch chan Message, quit chan bool, dec *gob.Decoder){
	for {
		select {
		case <- quit:
			return
		default:
			msg := new(Message)
			err := dec.Decode(msg)
			if err != nil {
				panic(err)
			}
			ch <- *msg
		}
	}
}
