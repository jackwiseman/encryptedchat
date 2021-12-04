package main

import (
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"net"
	"strings"
	"errors"
	"syscall"
	"io"
	"github.com/thanhpk/randstr"
)

type client struct {
	enc gob.Encoder
	dec gob.Decoder
	publicKey rsa.PublicKey
	isAuth bool
	authToken string
	conn net.Conn
	username string
	room *room
	commands chan<- command
}

type Message struct {
	Msg string
	Sender string
	PublicKey rsa.PublicKey
}

func (c *client) readInput() {

	for {
		msg := new(Message)
		err := c.dec.Decode(msg)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, io.EOF) || strings.Contains(err.Error(), "use of closed network connection"){
				return
			} else {
				panic(err)
			}
		}

		msgString := msg.Msg

		if msg.Sender == "cmd" {
			msgString, _ = decrypt(msg.Msg,SERVERKEY)
		}
		msgString = strings.Trim(msgString, "\r\n")
		args := strings.Split(msgString, " ")
		cmd := strings.TrimSpace(args[0])


		if(msg.Sender == "auth") {
			if(msg.Msg == c.authToken) {
				c.isAuth = true
				c.commands <- command{
					id:	CMD_AUTH,
					client: c,
					args: args,
				}
			}
			continue
		}

		if cmd == ""{
			c.publicKey = msg.PublicKey
			continue
		}

		switch cmd {
		case "/login":
			c.commands <- command{
				id:     CMD_LOGIN,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		case "/name":
			c.commands <- command{
				id: CMD_CHGNAME,
				client: c,
				args: args,
			}
		case "/help":
			c.commands <- command{
				id:	CMD_HELP,
				client: c,
				args:	args,
			}
		default:
			if cmd[0] == '/' {
				c.err(fmt.Errorf("Unknown command: %s, use /help to list all commands", cmd))
			} else {
				c.commands <- command{
					id:     CMD_MSG,
					client: c,
					args:   args,
				}
			}
		}
	}
}

func (c *client) err(err error) {
	message := new(Message)
	message.Msg = encrypt("ERROR: " + err.Error() + "\n", c.publicKey)
	message.Sender = "server"

	err2 := c.enc.Encode(message)
	if err2 != nil {
		panic(err2)
	}
}

func (c *client) msg(msg string, sender *client) {
	message := new(Message)
	message.Msg = msg + "\n"
	message.PublicKey = c.getEncryptionKey()
	message.Sender = sender.username

	err := c.enc.Encode(message)
	if err != nil {
		panic(err)
	}
}

func (c* client) eventMsg(msg string) {
	message := new(Message)
	message.Msg = encrypt("# " + msg + "\n", c.publicKey)
	message.Sender = "server"

	
	err := c.enc.Encode(message)
	if err != nil {
		panic(err)
	}
}

// update the client's public key for encryption
func (c *client) updateKey() {
	var key rsa.PublicKey
	for addr, k := range c.room.keys {
		if addr != c.conn.RemoteAddr() {
			key = k
		}
	}

	message := new(Message)
	message.PublicKey = key
	err := c.enc.Encode(message)
	if err != nil {
		panic(err)
	}
}

func (c *client) getEncryptionKey() rsa.PublicKey {
	var key rsa.PublicKey
	for addr, k := range c.room.keys {
		if addr != c.conn.RemoteAddr() {
			key = k
		}
	}
	return key
}

func (c *client) auth(key rsa.PublicKey) {
	c.authToken = randstr.String(16)
	token := encrypt(c.authToken, key)

	message := new(Message)
	message.PublicKey = key
	message.Msg = token
	message.Sender = "auth"
	err := c.enc.Encode(message)
	if err != nil {
		panic(err)
	}
}
