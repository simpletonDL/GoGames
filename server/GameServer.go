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
	// TODO: close connection?

	/* Client handler */
	clientManager := NewClientManager()
	go acceptClients(l, clientManager)

	/* Select team phase */

	selectTeamProcessor := NewGameProcessor(engine.SelectTeamMode, clientManager)
	go selectTeamProcessor.Run()
	<-selectTeamProcessor.ReadyToStart
	selectTeamProcessor.Cancel()

	/* Main game session */
	gameProcessor := NewGameProcessor(engine.MainGameMode, clientManager)
	for playerId, info := range selectTeamProcessor.GameEngine.Players {
		gameProcessor.GameEngine.Input <- engine.CreatePlayerCommand{
			Nickname: info.Nickname,
			Team:     info.Team,
			PlayerId: playerId,
			PosX:     settings.WorldWidth / 2,
			PosY:     settings.WorldHeight,
		}
	}
	go gameProcessor.Run()
	<-gameProcessor.Finished
}

func acceptClients(l net.Listener, manager *ClientManager) {
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

		manager.ConnectClient(client)
		manager.EnqueueCommand(engine.CreatePlayerCommand{
			Nickname: initCmd.Nickname,
			Team:     protocol.BlueTeam,
			PlayerId: engine.PlayerId(client.Id),
			PosX:     settings.WorldWidth / 2,
			PosY:     settings.WorldHeight,
		})
		currentClientId++
	}
}
