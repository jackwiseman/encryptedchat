package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Started server on :8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %s", err.Error())
			continue
		}

		// goroutine handler
		go s.newClient(conn)
	}
}
