package server

import (
	"encoding/json"
	"fmt"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"math/rand"
	"time"
)

/* Game step processing */

type GameProcessor struct {
	GameEngine *engine.GameEngine
	Clients    []Client
}

func NewGameProcessor() *GameProcessor {
	processor := &GameProcessor{
		GameEngine: engine.NewGameEngine(settings.GameInputCapacity),
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

func (p *GameProcessor) Run() {
	go AddRandomBoxEvent(p.GameEngine)
	p.GameEngine.Run(settings.GameFPS, settings.VelocityIterations, settings.PositionIterations)
}

func AddRandomBoxEvent(e *engine.GameEngine) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		e.ScheduleCommand(engine.NewCustomCommand(func(e *engine.GameEngine) {
			x := rand.Uint32() % settings.WorldWidth
			engine.AddBox(e.World, float64(x), settings.WorldHeight, 2, 1, 1, 0.5, 0.3)
		}))
	}
}
