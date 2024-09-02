package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/GoGames/common/utils"
	"log"
)

/* Game step processing */

type GameProcessor struct {
	GameEngine *engine.GameEngine
	Ctx        context.Context
	Cancel     context.CancelFunc

	Mod engine.GameEngineMod

	// SelectTeamMode
	ReadyToStart chan bool

	// MainGameMode
	Finished chan bool
}

func NewGameProcessor(mod engine.GameEngineMod, clientManager *ClientManager) *GameProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	var gameEngine *engine.GameEngine
	switch mod {
	case engine.SelectTeamMode:
		gameEngine = engine.NewSelectTeamGameEngine(ctx, mod, clientManager.Input)
	case engine.MainGameMode:
		gameEngine = engine.NewMainGameEngine(ctx, mod, clientManager.Input)
	default:
		log.Fatalf("Invalid game engine mode: $%d", mod)
	}

	processor := &GameProcessor{
		Ctx:          ctx,
		Cancel:       cancel,
		GameEngine:   gameEngine,
		Mod:          mod,
		ReadyToStart: make(chan bool, 1),
	}
	// This callback sends new game state to client every timestamp
	processor.GameEngine.AddListener(func(e *engine.GameEngine) {
		for _, client := range clientManager.GetAllClients() {
			encoder := json.NewEncoder(client.conn)
			state := engine.GetGameState(e)
			err := encoder.Encode(state)
			if err != nil {
				// TODO: remove client on disconnection
				fmt.Printf("Decodding game state error: %s\n", err.Error())
			}
		}
	})
	if mod == engine.SelectTeamMode {
		processor.GameEngine.AddListener(func(e *engine.GameEngine) {
			// Add change team handler
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

			// Add ready to start game handler
			everybodyIsReadyToStart := utils.AllEntries(e.Players, func(key engine.PlayerId, value *engine.PlayerInfo) bool {
				return value.IsReadyToStart
			})
			if everybodyIsReadyToStart && len(e.Players) > 0 {
				println("everybodyIsReadyToStart && len(e.Players) > 0")
				processor.ReadyToStart <- true
			}
		})
	}

	return processor
}

func (p *GameProcessor) Run() {

	p.GameEngine.Run(settings.GameFPS, settings.VelocityIterations, settings.PositionIterations)
}
