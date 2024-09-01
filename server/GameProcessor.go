package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
)

/* Game step processing */

type GameProcessor struct {
	GameEngine *engine.GameEngine
	Clients    []Client

	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewGameProcessor() *GameProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	processor := &GameProcessor{
		Ctx:        ctx,
		Cancel:     cancel,
		GameEngine: engine.NewGameEngine(ctx, settings.GameInputCapacity),
		Clients:    []Client{},
	}
	// This callback sends new game state to client every timestamp
	processor.GameEngine.AddListener(func(gameState protocol.GameState) {
		for _, client := range processor.Clients {
			encoder := json.NewEncoder(client.conn)
			err := encoder.Encode(gameState)
			if err != nil {
				// TODO: remove client on disconnection
				fmt.Printf("Decodding game state error: %s\n", err.Error())
			}
		}
	})
	return processor
}

func (p *GameProcessor) ConnectClient(client Client) {
	p.Clients = append(p.Clients, client)
	go HandleClientInput(client, p)
}

func (p *GameProcessor) Run() {
	p.GameEngine.Run(settings.GameFPS, settings.VelocityIterations, settings.PositionIterations)
}
