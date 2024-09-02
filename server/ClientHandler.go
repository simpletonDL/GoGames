package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/GoGames/common/utils"
	"net"
	"sync"
)

type Client struct {
	Id       engine.PlayerId
	Nickname string
	conn     net.Conn
	Decoder  *json.Decoder
	Encoder  *json.Encoder
}

func NewClient(id uint8, conn net.Conn) Client {
	return Client{
		Id:       engine.PlayerId(id),
		Nickname: "",
		conn:     conn,
		Decoder:  json.NewDecoder(conn),
		Encoder:  json.NewEncoder(conn),
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

func (h *ClientManager) addClient(client Client) {
	h.clients = append(h.clients, client)
}

func (h *ClientManager) AddClient(client Client) {
	h.mu.Lock()
	h.addClient(client)
	h.mu.Unlock()
}

func (h *ClientManager) removeClient(id engine.PlayerId) {
	h.clients = utils.Filter(h.clients, func(client Client) bool { return client.Id != id })
}

func (h *ClientManager) RemoveClient(id engine.PlayerId) {
	h.mu.Lock()
	h.removeClient(id)
	h.mu.Unlock()
}

func (h *ClientManager) GetAllClients() []Client {
	var clients []Client
	h.mu.Lock()
	clients = h.clients
	h.mu.Unlock()
	return clients
}

func (h *ClientManager) broadcast(state protocol.GameState) {
	for _, client := range h.clients {
		err := client.Encoder.Encode(state)
		if err != nil {
			h.removeClient(client.Id)
			fmt.Printf("Disconect client %s with error %s\n", client.Nickname, err.Error())
		}
	}
}

func (h *ClientManager) Broadcast(state protocol.GameState) {
	h.mu.Lock()
	h.broadcast(state)
	h.mu.Unlock()
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
			h.RemoveClient(client.Id)
			fmt.Printf("Client %s disconected with err %s\n", client.Nickname, err.Error())
			break
		}
		h.EnqueueCommand(engine.PlayerInputCommand{PlayerId: client.Id, Cmd: cmd})
	}
}
