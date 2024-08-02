package server

import (
	"encoding/json"
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"io"
	"net"
)

type Client struct {
	Id   uint8
	conn net.Conn
}

func HandleClientInput(client Client, processor *GameProcessor) {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	for {
		var cmd protocol.ClientInputCommand
		if err := decoder.Decode(&cmd); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error decoding JSON: %s\n", err.Error())
			continue
		}
		processor.GameEngine.ScheduleCommand(engine.PlayerInputCommand{PlayerId: engine.PlayerId(client.Id), Cmd: cmd})
	}
	fmt.Printf("Connection closed: %s\n", client.conn.RemoteAddr())
}
