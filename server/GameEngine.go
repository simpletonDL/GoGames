package server

import (
	"encoding/json"
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/simpletonDL/GoGames/common/engine"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"net"
	"time"
)

/* Initial world */

func createInitialWorld() *box2d.B2World {
	world := engine.NewWorld(0, -10)
	// Ground body
	engine.AddBox(world, box2d.B2BodyType.B2_staticBody, 8, 1, 0, 16, 2, 0, 1)
	// Dynamic body
	engine.AddBox(world, box2d.B2BodyType.B2_dynamicBody, 8, 15, 1, 1, 1, 1, 0.3)
	return world
}

/* Game step processing */

type GameEngine struct {
	input   chan protocol.InputCommand
	world   *box2d.B2World
	clients []net.Conn // TODO: refactor
}

func NewGameEngine() *GameEngine {
	return &GameEngine{
		world:   createInitialWorld(),
		input:   make(chan protocol.InputCommand, settings.ServerInputCapacity),
		clients: []net.Conn{},
	}
}

func (p *GameEngine) Run() {
	timestamp := 0
	ticker := time.NewTicker(time.Second / settings.ServerFPS)
	for {
		<-ticker.C
		commands := p.collectAllInputCommandsNonBlocking()
		p.processInputCommands(commands)
		p.world.Step(settings.ServerTimestamp, settings.VelocityIterations, settings.PositionIterations)
		p.notifyClients()
		timestamp++
	}
}

func (p *GameEngine) collectAllInputCommandsNonBlocking() []protocol.InputCommand {
	result := []protocol.InputCommand{}
	for {
		select {
		case input := <-p.input:
			result = append(result, input)
		default:
			return result
		}
	}
}

func (p *GameEngine) processInputCommands(commands []protocol.InputCommand) {
	for _, cmd := range commands {
		switch cmd.Id {
		case protocol.MouseClickID:
			x := cmd.Args["x"]
			y := cmd.Args["y"]
			fmt.Printf("processInputCommands: Click %f %f\n", x, y)
			engine.AddBox(p.world, box2d.B2BodyType.B2_dynamicBody, x, y, 1, 1, 1, 1, 0.3)
		}
	}
}

func (p *GameEngine) notifyClients() {
	for _, client := range p.clients {
		encoder := json.NewEncoder(client)
		gameState := engine.B2WorldToGameState(p.world)
		err := encoder.Encode(gameState)
		if err != nil {
			// TODO: remove client on disconnection
			fmt.Printf("Decodding game state error: %s\n", err.Error())
		}
	}
}

func (p *GameEngine) SendCommand(command protocol.InputCommand) {
	p.input <- command
}
