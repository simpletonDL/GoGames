package server

import (
	"encoding/json"
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/settings"
)

/* Game step processing */

type GameProcessor struct {
	GameEngine *engine.GameEngine
	Clients    []Client
}

func NewGameProcessor() *GameProcessor {
	processor := &GameProcessor{
		GameEngine: engine.NewGameEngine(settings.ServerInputCapacity),
		Clients:    []Client{},
	}
	// This callback sends new game state to client every timestamp
	processor.GameEngine.AddListener(func(world *box2d.B2World) {
		for _, client := range processor.Clients {
			encoder := json.NewEncoder(client.conn)
			gameState := engine.B2WorldToGameState(world)
			err := encoder.Encode(gameState)
			if err != nil {
				// TODO: remove client on disconnection
				fmt.Printf("Decodding game state error: %s\n", err.Error())
			}
		}
	})
	return processor
}

func (p *GameProcessor) Run() {
	p.GameEngine.Run(settings.ServerFPS, settings.VelocityIterations, settings.PositionIterations)
}