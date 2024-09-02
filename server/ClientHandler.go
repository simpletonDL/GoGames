package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/GoGames/common/utils"
	"io"
	"net"
	"sync"
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

type ClientManager struct {
	Ctx     context.Context
	Cancel  context.CancelFunc
	Input   chan engine.GameCommand
	clients []Client
	mu      sync.Mutex
}

func NewClientManager() *ClientManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClientManager{
		Ctx:    ctx,
		Cancel: cancel,
		Input:  make(chan engine.GameCommand, settings.GameInputCapacity),
		mu:     sync.Mutex{},
	}
}

func (h *ClientManager) EnqueueCommand(cmd engine.GameCommand) {
	h.Input <- cmd
}

func (h *ClientManager) AddClient(client Client) {
	h.mu.Lock()
	h.clients = append(h.clients, client)
	h.mu.Unlock()
}

func (h *ClientManager) GetAllClients() []Client {
	var clients []Client
	h.mu.Lock()
	clients = h.clients
	h.mu.Unlock()
	return clients
}

func (h *ClientManager) ConnectClient(client Client) {
	h.AddClient(client)
	go h.HandleClientInput(client)
}

func (h *ClientManager) HandleClientInput(client Client) {
	defer client.conn.Close()
	for {
		select {
		case <-h.Ctx.Done():
			utils.Log("Client handler cancelled (id=%d)\n", client.Id)
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
		h.EnqueueCommand(engine.PlayerInputCommand{PlayerId: engine.PlayerId(client.Id), Cmd: cmd})
	}
	fmt.Printf("Connection closed: %s\n", client.conn.RemoteAddr())
}
