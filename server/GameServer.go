package server

import (
	"fmt"
	"net"
)

func Run(port string) {
	l, _ := net.Listen("tcp4", port)
	defer l.Close()

	engine := NewGameEngine()
	go engine.Run()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Connection error: %s\n", err.Error())
			continue
		}
		fmt.Printf("New connection from %s\n", conn.RemoteAddr())
		engine.clients = append(engine.clients, conn)
		go HandleClientInput(conn, engine)
	}
}
