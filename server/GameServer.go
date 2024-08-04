package server

import (
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"net"
)

func Run(port string) {
	l, _ := net.Listen("tcp4", port)
	defer l.Close()

	processor := NewGameProcessor()
	go processor.Run()
	currentClientId := uint8(0)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Connection error: %s\n", err.Error())
			continue
		}
		fmt.Printf("New connection from %s\n", conn.RemoteAddr())
		client := Client{
			Id:   currentClientId,
			conn: conn,
		}
		processor.Clients = append(processor.Clients, client)
		processor.GameEngine.ScheduleCommand(engine.CreatePlayerCommand{
			PlayerId: engine.PlayerId(client.Id),
			PosX:     2,
			PosY:     15,
		})
		go HandleClientInput(client, processor)
		currentClientId++
	}
}
