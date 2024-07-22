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
	if len(os.Args) < 2 {
		log.Fatal("Port number is required as a program argument")
	}

	port := os.Args[1]
	address := "localhost:" + port

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
