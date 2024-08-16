package engine

import (
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/box2d"
	"time"
)

type GameEngine struct {
	Input     chan GameCommand
	World     *box2d.B2World
	Players   map[PlayerId]PlayerInfo
	Listeners []GameEngineListener
}

func NewGameEngine(inputCapacity int) *GameEngine {
	engine := &GameEngine{
		World:     createInitialWorld(),
		Input:     make(chan GameCommand, inputCapacity),
		Players:   map[PlayerId]PlayerInfo{},
		Listeners: []GameEngineListener{},
	}
	// Add collision logic
	engine.World.SetContactListener(NewCollisionTracker(engine))
	return engine
}

func (e *GameEngine) Run(fps int, velocityIterations int, positionIterations int) {
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	timestamp := 1.0 / float64(fps)
	for {
		<-ticker.C
		commands := e.collectAllGameCommandsNonBlocking()
		for _, command := range commands {
			command.Execute(e)
		}
		e.World.Step(timestamp, velocityIterations, positionIterations)
		for _, listener := range e.Listeners {
			listener(e.World)
		}
	}
}

func (e *GameEngine) ScheduleCommand(cmd GameCommand) {
	e.Input <- cmd
}

type PlayerId uint8

var PlayerDirection = struct {
	Left  bool
	Right bool
}{
	Left:  false,
	Right: true,
}

type PlayerInfo struct {
	Body      *box2d.B2Body
	Direction bool // right=true, left=false
}

type GameEngineListener func(world *box2d.B2World)

func (e *GameEngine) AddListener(listener GameEngineListener) {
	e.Listeners = append(e.Listeners, listener)
}

func (e *GameEngine) collectAllGameCommandsNonBlocking() []GameCommand {
	var result []GameCommand
	for {
		select {
		case input := <-e.Input:
			result = append(result, input)
		default:
			return result
		}
	}
}

func B2WorldToGameState(world *box2d.B2World) protocol.GameState {
	gameObjects := make([]protocol.GameObject, world.GetBodyCount())
	for body := world.GetBodyList(); body != nil; body = body.M_next {
		data := body.GetUserData().(BodyUserData)
		object := protocol.GameObject{
			XPos:      body.GetPosition().X,
			YPos:      body.GetPosition().Y,
			Angel:     body.GetAngle(),
			ImageKind: data.Kind,
			Width:     data.Width,
			Height:    data.Height,
		}
		gameObjects = append(gameObjects, object)
	}
	return protocol.GameState{
		Objects: gameObjects,
	}
}
