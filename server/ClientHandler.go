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
	Id      uint8
	conn    net.Conn
	Decoder *json.Decoder
}

func NewClient(id uint8, conn net.Conn) Client {
	return Client{
		Id:      id,
		conn:    conn,
		Decoder: json.NewDecoder(conn),
	}
}

func Receive[T any](client Client) (result T, err error) {
	err = client.Decoder.Decode(&result)
	return
}

func (p *GameProcessor) HandleClientInput(client Client) {
	defer client.conn.Close()
	for {
		select {
		case <-p.Ctx.Done():
			return
		default:
		}
		cmd, err := Receive[protocol.ClientInputCommand](client)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error decoding JSON: %s\n", err.Error())
			continue
		}
		p.GameEngine.ScheduleCommand(engine.PlayerInputCommand{PlayerId: engine.PlayerId(client.Id), Cmd: cmd})
	}
	fmt.Printf("Connection closed: %s\n", client.conn.RemoteAddr())
}
