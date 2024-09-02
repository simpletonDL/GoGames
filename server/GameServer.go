package server

import (
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"net"
)

func Run(port string) {
	l, _ := net.Listen("tcp4", port)
	defer l.Close()

	processor := NewGameProcessor(engine.SelectTeamMode)
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
			Team:     protocol.BlueTeam,
			PlayerId: engine.PlayerId(client.Id),
			PosX:     settings.WorldWidth / 2,
			PosY:     settings.WorldHeight,
		})
		currentClientId++
	}
}
