package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	username     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

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
	c.conn.Write([]byte("ERROR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte(msg + "\n"))
}

func (c* client) eventMsg(msg string) {
	c.conn.Write([]byte("# " + msg + "\n"))
}
