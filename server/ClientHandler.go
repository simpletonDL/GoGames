package server

import (
	"encoding/json"
	"fmt"
	"github.com/simpletonDL/GoGames/common/protocol"
	"io"
	"net"
)

func HandleClientInput(conn net.Conn, engine *GameEngine) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	for {
		var cmd protocol.InputCommand
		if err := decoder.Decode(&cmd); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error decoding JSON: %s\n", err.Error())
			continue
		}
		engine.SendCommand(cmd)
	}
	fmt.Printf("Connection closed: %s\n", conn.RemoteAddr())
}
