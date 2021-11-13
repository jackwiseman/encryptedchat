package main

import "net"

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string, isEvent bool) {
	for addr, m := range r.members {
		if addr != sender.conn.RemoteAddr() {
			if isEvent {
				m.eventMsg(msg)
			} else {
				m.msg(msg)
			}
		}
	}
}
