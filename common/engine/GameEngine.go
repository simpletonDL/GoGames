package engine

import (
	"context"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/GoGames/common/utils"
	"github.com/simpletonDL/box2d"
	"math/rand"
	"time"
)

type GameEngine struct {
	Ctx       context.Context
	Input     chan GameCommand
	World     *box2d.B2World
	Players   map[PlayerId]*PlayerInfo
	Listeners []GameEngineListener
	Events    []GameEvent
	Mod       GameEngineMod
}

type GameEngineMod uint8

const (
	SelectTeamMode = GameEngineMod(iota)
	MainGameMode
)

func NewGameEngine(ctx context.Context, createWorldFun func() *box2d.B2World, events []GameEvent, mod GameEngineMod, InputQueue chan GameCommand) *GameEngine {
	engine := &GameEngine{
		Ctx:       ctx,
		World:     createWorldFun(),
		Input:     InputQueue,
		Players:   map[PlayerId]*PlayerInfo{},
		Listeners: []GameEngineListener{},
		Events:    events,
		Mod:       mod,
	}
	// Add collision logic
	engine.World.SetContactListener(NewCollisionTracker(engine))
	return engine
}

func NewMainGameEngine(ctx context.Context, mod GameEngineMod, inputQueue chan GameCommand) *GameEngine {
	events := []GameEvent{
		NewWeaponBoxCreationEvent(time.Second*10, 2),
		NewBoxCreationEvent(time.Second*10, 2),
	}
	return NewGameEngine(ctx, createMainGameWorld, events, mod, inputQueue)
}

func NewSelectTeamGameEngine(ctx context.Context, mod GameEngineMod, inputQueue chan GameCommand) *GameEngine {
	return NewGameEngine(ctx, createSelectTeamWorld, []GameEvent{}, mod, inputQueue)
}

func (e *GameEngine) Run(fps int, velocityIterations int, positionIterations int) {
	// Run all events async
	for _, event := range e.Events {
		go e.RunEvent(event)
	}

	ticker := time.NewTicker(time.Second / time.Duration(fps))
	timestamp := 1.0 / float64(fps)
	for {
		<-ticker.C
		select {
		case <-e.Ctx.Done():
			utils.Log("GameEngine Box2D loop cancelled\n")
			return
		default:
		}
		commands := e.collectAllGameCommandsNonBlocking()
		for _, command := range commands {
			command.Execute(e)
		}
		e.World.Step(timestamp, velocityIterations, positionIterations)
		e.processOutOfScreenBodies()
		e.processPlayersDeath()
		e.processWeaponsReload()
		for _, listener := range e.Listeners {
			listener(e)
		}
	}
}

func (e *GameEngine) ScheduleCommand(cmd GameCommand) {
	e.Input <- cmd
}

type PlayerId uint8

type PlayerInfo struct {
	Nickname              string
	Team                  protocol.TeamKind
	Body                  *box2d.B2Body
	Direction             protocol.DirectionKind
	MoveDownThrowPlatform bool
	JumpCount             int8
	Weapon                Weapon
	LivesCount            int8
	IsAlive               bool

	// SelectTeamMode
	IsReadyToStart bool
}

type GameEngineListener func(engine *GameEngine)

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

		object := protocol.GameObject{
			XPos:     body.GetPosition().X,
			YPos:     body.GetPosition().Y,
			Angel:    body.GetAngle(),
			BodyKind: data.GetKind(),
			Width:    data.GetWidth(),
			Height:   data.GetHeight(),
		}
		if data.GetKind() == protocol.BodyKindHero {
			playerInfo := engine.Players[data.(PlayerUserData).HeroId]
			weaponInfo := playerInfo.Weapon.GetInfo()
			object.Nickname = playerInfo.Nickname
			object.Team = playerInfo.Team
			object.Direction = playerInfo.Direction
			object.WeaponKind = weaponInfo.WeaponKind
			object.WeaponAvailableBullets = weaponInfo.WeaponAvailableBullets
			object.WeaponAvailableBulletsInMagazine = weaponInfo.WeaponAvailableBulletsInMagazine
			object.WeaponMagazineCapacity = weaponInfo.WeaponMagazineCapacity
			object.WeaponIsReady = weaponInfo.WeaponIsReady
			object.LivesCount = int(playerInfo.LivesCount)
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

				playerInfo := e.Players[userData.(PlayerUserData).HeroId]
				playerInfo.LivesCount -= 1
			} else {
				// remove body
				e.World.DestroyBody(body)
			}
		}
	}
}

func (e *GameEngine) processPlayersDeath() {
	for _, playerInfo := range e.Players {
		if playerInfo.LivesCount == 0 && playerInfo.IsAlive {
			playerInfo.IsAlive = false
			e.World.DestroyBody(playerInfo.Body)
		}
	}
}

func (e *GameEngine) processWeaponsReload() {
	for _, info := range e.Players {
		info.Weapon.ProcessGameTick()
	}
}
