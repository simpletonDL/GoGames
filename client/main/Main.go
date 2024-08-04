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
	client.LoadImages()
}

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Required arguments: host port nickname (like `localhost 5005 rich.bitch`)")
	}

	host := os.Args[1]
	port := os.Args[2]
	nickname := os.Args[3]
	_ = nickname

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
