package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/simpletonDL/GoGames/client"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"log"
	"net"
	"os"
)

func init() {
	client.LoadImages("./assets")
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Host and port is required (like localhost 5005)")
	}
	if argsCount := len(os.Args); argsCount > 2 {

	}

	host := os.Args[1]
	port := os.Args[2]
	address := host + ":" + port

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	game := &client.Game{
		Conn:      conn,
		GameState: protocol.NewEmptyGameState(),
	}
	go game.HandleServerGameState()

	ebiten.SetWindowSize(settings.ScreenWidth, settings.ScreenHeight)
	if err := ebiten.RunGame(game); err == nil {
		log.Fatal(err)
	}
}
