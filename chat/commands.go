package main

type commandID int

// All possible commands
const (
	CMD_LOGIN commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_HELP
)

type command struct {
	id     commandID
	client *client
	args   []string
}
