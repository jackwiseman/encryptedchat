package main

import (
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"net"
	"strings"
)

type client struct {
	enc      gob.Encoder
	dec      gob.Decoder
	conn     net.Conn
	username     string
	room     *room
	commands chan<- command
}

type Message struct {
	Msg string
	Publickey rsa.PublicKey
}

func (c *client) readInput() {

	for {
		msg := new(Message)
		err := c.dec.Decode(msg)
		if err != nil {
			panic(err)
		}
		msgString := msg.Msg
		msgString = strings.Trim(msgString, "\r\n")

		args := strings.Split(msgString, " ")
		cmd := strings.TrimSpace(args[0])

		if cmd == ""{
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
	message.Msg = "ERROR: " + err.Error() + "\n"

	err2 := c.enc.Encode(message)
	if err2 != nil {
		panic(err2)
	}
}

func (c *client) msg(msg string) {
	message := new(Message)
	message.Msg = msg + "\n"

	err := c.enc.Encode(message)
	if err != nil {
		panic(err)
	}
}

func (c* client) eventMsg(msg string) {
	message := new(Message)
	message.Msg = "# " + msg + "\n"

	err := c.enc.Encode(message)
	if err != nil {
		panic(err)
	}
}
