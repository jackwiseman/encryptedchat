package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_LOGIN:
			s.username(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		case CMD_HELP:
			s.help(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New Client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		username:     "Anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) username(c *client, args []string) {
	c.username = args[1]
	c.msg(fmt.Sprintf("You are logged in as %s", c.username))
}

func (s *server) join(c *client, args []string) {

	if len(args) != 2 {
		c.msg(fmt.Sprintf("Please specify a room name you would like to join"))
		return
	}

	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r
	r.broadcast(c, fmt.Sprintf("%s has joined the room.", c.username))
	c.msg(fmt.Sprintf("Welcome to the room %s", r.name))
}

func (s *server) listRooms(c *client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	if len(rooms) == 0 {
		c.msg(fmt.Sprintf("There are no active rooms, create one with /join {name}"))
	} else {
		c.msg(fmt.Sprintf("Available rooms: %s", strings.Join(rooms, ", ")))
	}
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("You must join a room before sending a message"))
		return
	}

	c.room.broadcast(c, c.username+": "+strings.Join(args, " "))
}

func (s *server) quit(c *client, args []string) {
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)
	c.msg("Come back soon!")
	c.conn.Close()
}

func (s *server) help(c *client, args[]string) {
	c.msg("/join {room name}, /rooms, /quit")
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room.", c.username))
	}
}
