package main

import (
	"encoding/gob"
//	"encoding/hex"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"strings"
)

var SERVERKEY rsa.PrivateKey

type server struct {
	rooms    map[string]*room
	commands chan command
	users map[string]rsa.PublicKey
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
		users:    make(map[string]rsa.PublicKey),
	}
}

func (s *server) run() {
	loadPrivateKey(&SERVERKEY)
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_LOGIN:
			s.login(cmd.client, cmd.args)
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
		case CMD_AUTH:
			s.authenticated(cmd.client)
		case CMD_CHGNAME:
			s.changeName(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New Client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		enc: *gob.NewEncoder(conn),
		dec: *gob.NewDecoder(conn),
		isAuth: false,
		conn:     conn,
		username:     "",
		commands: s.commands,
	}

	s.sendPublicKeyToClient(c)
	c.readInput()
	var args []string
	s.quit(c, args)
}

func (s *server) login(c *client, args []string) {
	if c.isAuth {
		c.eventMsg("You are already logged in")
		return
	}

	if len(args) != 2 {
		c.eventMsg("Your username must be one word, please try again")
		return
	}

	if args[1] == "server" || args[1] == "serverkey" || args[1] == "auth" || args[1] == "keyX" || len(args[1]) == 0 {
		c.eventMsg("Illegal name choice, please try again")
		return
	}

	for user, k := range s.users {
		if k.E == c.publicKey.E && k.N.Cmp(c.publicKey.N) == 0 {
			if args[1] != user {
				c.eventMsg("You cannot login as that user, please try again")
				return
			}
		}
	}

	c.username = args[1]

	key, ok := s.users[c.username]
	if ok { // need to auth
		c.eventMsg(fmt.Sprintf("Logging in as %s...", c.username))
		c.auth(key)
	} else { // new account, no need to auth
		s.users[c.username] = c.publicKey
		c.isAuth = true
		c.eventMsg(fmt.Sprintf("Welcome new user %s, your connection is now secure", c.username))
	}
}

func (s *server) changeName(c *client, args []string) {
	if !c.isAuth {
		c.eventMsg("You must be logged in to change your name")
		return
	}
	if len(args) != 2 {
		c.eventMsg("Your username must be one word, please try again")
		return
	}
	for user, _ := range s.users {
		if args[1] == user {
			c.eventMsg("That username is already taken, please try again")
			return
		}
	}
	if args[1] == "server" || args[1] == "serverkey" || args[1] == "auth" || args[1] == "keyX" || len(args[1]) == 0 {
		c.eventMsg("Illegal name choice, please try again")
		return
	}

	c.username = args[1]
	delete(s.users, c.username)
	s.users[c.username] = c.publicKey
	c.eventMsg(fmt.Sprintf("Successfully changed username to %s", c.username))
}

func (s *server) authenticated(c *client) {
	c.eventMsg(fmt.Sprintf("Sucess: You are logged in as %s", c.username))
}

func (s *server) join(c *client, args []string) {

	if !c.isAuth {
		c.eventMsg("You must login before joining a room!")
		return
	}

	if len(args) != 2 || len(args[1]) == 0 {
		c.eventMsg("Please specify a room name you would like to join (rooms must be one word)")
		return
	}

	roomName := args[1]
	r, exists := s.rooms[roomName]
	if !exists {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
			keys: make(map[net.Addr]rsa.PublicKey),
		}
		s.rooms[roomName] = r
	}

	_, ok := r.members[c.conn.RemoteAddr()]
	if ok {
		c.eventMsg("You are already in this room")
		return
	}

	if len(r.members) >= 2 {
		c.eventMsg("That room is full, try a different one or create your own")
		return
	} else {
		r.members[c.conn.RemoteAddr()] = c
		r.keys[c.conn.RemoteAddr()] = c.publicKey
	}

	s.quitCurrentRoom(c)
	c.room = r
	r.broadcast(c, fmt.Sprintf("%s has joined the room.", c.username), true)
	c.eventMsg(fmt.Sprintf("Welcome to the room %s", r.name))
}

func (s *server) listRooms(c *client, args []string) {
	if !c.isAuth {
		c.eventMsg("You must login before using this command")
		return
	}

	var rooms []string
	for name := range s.rooms {
		if len(s.rooms[name].members) < 2 {
			rooms = append(rooms, name)
		}
	}

	if len(rooms) == 0 {
		c.eventMsg("There are no joinable rooms, create one with /join {name}")
	} else {
		c.eventMsg(fmt.Sprintf("Available rooms: %s", strings.Join(rooms, ", ")))
	}
}

func (s *server) msg(c *client, args []string) {
	if !c.isAuth {
		c.eventMsg("You must login before sending a message - /login {username}")
		return
	}
	if c.room == nil {
		c.eventMsg("You must join a room before sending a message - /join {room name}")
		return
	}

	c.room.broadcast(c, strings.Join(args, " "), false)
}

func (s *server) quit(c *client, args []string) {
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)
	c.conn.Close()
}

func (s *server) help(c *client, args[]string) {
	c.eventMsg("/join {room name}, /rooms, /quit, /name {new username}")
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		if len(c.room.members) > 1 {
			delete(c.room.members, c.conn.RemoteAddr())
			delete(c.room.keys, c.conn.RemoteAddr())
			c.room.broadcast(c, fmt.Sprintf("%s has left the room.", c.username), true)
		} else {
			delete(s.rooms, c.room.name)
		}
	}
}

func (s *server) sendPublicKeyToClient(c *client) {
	msg := new(Message)
	msg.PublicKey = SERVERKEY.PublicKey
	msg.Sender = "serverkey"
	err := c.enc.Encode(msg)
	if err != nil {
		panic(err)
	}
}
