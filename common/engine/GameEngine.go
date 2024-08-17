package engine

import (
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/box2d"
	"math/rand"
	"time"
)

type GameEngine struct {
	Input     chan GameCommand
	World     *box2d.B2World
	Players   map[PlayerId]*PlayerInfo
	Listeners []GameEngineListener
}

func NewGameEngine(inputCapacity int) *GameEngine {
	engine := &GameEngine{
		World:     createInitialWorld(),
		Input:     make(chan GameCommand, inputCapacity),
		Players:   map[PlayerId]*PlayerInfo{},
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
		e.processOutOfScreenBodies()
		for _, listener := range e.Listeners {
			listener(GetGameState(e))
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
	Body                  *box2d.B2Body
	Direction             bool // right=true, left=false
	MoveDownThrowPlatform bool
	JumpCount             int8
}

type GameEngineListener func(world protocol.GameState)

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

func GetGameState(engine *GameEngine) protocol.GameState {
	world := engine.World
	gameObjects := make([]protocol.GameObject, world.GetBodyCount())
	for body := world.GetBodyList(); body != nil; body = body.M_next {
		data := GetBodyUserData(body)
		direction := true
		if data.GetKind() == protocol.BodyKindHero {
			direction = engine.Players[data.(PlayerUserData).HeroId].Direction
		}

		object := protocol.GameObject{
			XPos:      body.GetPosition().X,
			YPos:      body.GetPosition().Y,
			Angel:     body.GetAngle(),
			BodyKind:  data.GetKind(),
			Width:     data.GetWidth(),
			Height:    data.GetHeight(),
			Direction: direction,
		}
		gameObjects = append(gameObjects, object)
	}
	return protocol.GameState{
		Objects: gameObjects,
	}
}

func (e *GameEngine) processOutOfScreenBodies() {
	for body := e.World.GetBodyList(); body != nil; body = body.M_next {
		bodyIsOutOfBound := false
		x := body.GetPosition().X
		y := body.GetPosition().Y
		if x < 0-settings.OutOfScreenBound || x > settings.WorldWidth+settings.OutOfScreenBound {
			bodyIsOutOfBound = true
		}
		if y < 0-settings.OutOfScreenBound || y > settings.WorldHeight+settings.OutOfScreenBound {
			bodyIsOutOfBound = true
		}
		if bodyIsOutOfBound {
			userData := GetBodyUserData(body)
			if userData.GetKind() == protocol.BodyKindHero {
				// respawn hero, TODO: decrement life count
				newX := rand.Uint32() % settings.WorldWidth
				newY := settings.WorldHeight
				body.SetTransform(box2d.B2Vec2{X: float64(newX), Y: float64(newY)}, 0)
				body.SetLinearVelocity(box2d.B2Vec2{X: 0, Y: 0})
			} else {
				// remove body
				e.World.DestroyBody(body)
			}
		}
	}
}
