package server

import (
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
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
		client := NewClient(currentClientId, conn)

		initCmd, err := Receive[protocol.ClientInitializationCommand](client)
		if err != nil {
			fmt.Printf("Error during client initializtion: %s\n", err.Error())
			continue
		}

		processor.ConnectClient(client)
		processor.GameEngine.ScheduleCommand(engine.CreatePlayerCommand{
			Nickname: initCmd.Nickname,
			PlayerId: engine.PlayerId(client.Id),
			PosX:     2,
			PosY:     15,
		})
		currentClientId++
	}
}
