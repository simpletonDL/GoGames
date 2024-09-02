package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"log"
)

/* Game step processing */

type GameProcessor struct {
	GameEngine *engine.GameEngine
	Clients    []Client

	Ctx    context.Context
	Cancel context.CancelFunc

	Mod engine.GameEngineMod
}

func NewGameProcessor(mod engine.GameEngineMod) *GameProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	var gameEngine *engine.GameEngine
	switch mod {
	case engine.SelectTeamMode:
		gameEngine = engine.NewSelectTeamGameEngine(ctx, mod)
	case engine.MainGameMode:
		gameEngine = engine.NewMainGameEngine(ctx, mod)
	default:
		log.Fatalf("Invalid game engine mode: $%d", mod)
	}

	processor := &GameProcessor{
		Ctx:        ctx,
		Cancel:     cancel,
		GameEngine: gameEngine,
		Clients:    []Client{},
		Mod:        mod,
	}
	// This callback sends new game state to client every timestamp
	processor.GameEngine.AddListener(func(e *engine.GameEngine) {
		for _, client := range processor.Clients {
			encoder := json.NewEncoder(client.conn)
			state := engine.GetGameState(e)
			err := encoder.Encode(state)
			if err != nil {
				// TODO: remove client on disconnection
				fmt.Printf("Decodding game state error: %s\n", err.Error())
			}
		}
	})
	// Add change team handler
	if mod == engine.SelectTeamMode {
		processor.GameEngine.AddListener(func(e *engine.GameEngine) {
			for _, playerInfo := range e.Players {
				xPos := playerInfo.Body.GetPosition().X
				var playerTeam protocol.TeamKind
				if xPos > settings.WorldWidth/2 {
					playerTeam = protocol.RedTeam
				} else {
					playerTeam = protocol.BlueTeam
				}
				playerInfo.Team = playerTeam
			}
		})
	}

	return processor
}

func (p *GameProcessor) ConnectClient(client Client) {
	p.Clients = append(p.Clients, client)
	go p.HandleClientInput(client)
}

func (p *GameProcessor) Run() {
	p.GameEngine.Run(settings.GameFPS, settings.VelocityIterations, settings.PositionIterations)
}
