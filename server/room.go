package main

import (
	"net"
	"crypto/rsa"
)

type room struct {
	name    string
	members map[net.Addr]*client
	keys map[net.Addr]rsa.PublicKey
}

func (r *room) broadcast(sender *client, msg string, isEvent bool) {
	for addr, m := range r.members {
		if addr != sender.conn.RemoteAddr() {
			if isEvent {
				m.eventMsg(msg)
				m.updateKey()
			} else {
				m.msg(msg, sender)
			}
		} else {
			m.updateKey()
		}
	}
}
