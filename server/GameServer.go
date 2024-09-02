package server

import (
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/GoGames/common/utils"
	"net"
)

func Run(port string) {
	l, _ := net.Listen("tcp4", port)
	// TODO: close connection?

	/* Client handler */
	clientManager := NewClientManager()
	go acceptClients(l, clientManager)

	/* Select team phase */

	for {
		selectTeamProcessor := NewGameProcessor(engine.SelectTeamMode, clientManager)
		// TODO: make it thread-safe
		currentClients := clientManager.GetAllClients()
		for _, client := range currentClients {
			addPlayerToWorld(selectTeamProcessor, client.Nickname, protocol.BlueTeam, client.Id)
		}

		go selectTeamProcessor.Run()
		<-selectTeamProcessor.ReadyToStart
		selectTeamProcessor.Cancel()

		/* Main game session */
		gameProcessor := NewGameProcessor(engine.MainGameMode, clientManager)
		for playerId, info := range selectTeamProcessor.GameEngine.Players {
			addPlayerToWorld(gameProcessor, info.Nickname, info.Team, playerId)
		}
		go gameProcessor.Run()
		winner := <-gameProcessor.TeamWin
		utils.Log("Winner: %s", winner.ToString())
		gameProcessor.Cancel()
	}
}

func addPlayerToWorld(gameProcessor *GameProcessor, nickname string, team protocol.TeamKind, playerId engine.PlayerId) {
	gameProcessor.GameEngine.Input <- engine.CreatePlayerCommand{
		Nickname:   nickname,
		Team:       team,
		PlayerId:   playerId,
		PosX:       settings.WorldWidth / 2,
		PosY:       settings.WorldHeight,
		LivesCount: settings.PlayerLivesCount,
	}
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
		client.Nickname = initCmd.Nickname

		manager.ConnectClient(client)
		manager.EnqueueCommand(engine.CreatePlayerCommand{
			Nickname:   initCmd.Nickname,
			Team:       protocol.BlueTeam,
			PlayerId:   client.Id,
			PosX:       settings.WorldWidth / 2,
			PosY:       settings.WorldHeight,
			LivesCount: 100,
		})
		currentClientId++
	}
}
