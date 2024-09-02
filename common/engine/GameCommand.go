package engine

import (
	"fmt"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
	"github.com/simpletonDL/box2d"
)

type GameCommand interface {
	Execute(engine *GameEngine)
}

type CustomCommand struct {
	ExecuteFun func(engine *GameEngine)
}

func (c CustomCommand) Execute(engine *GameEngine) {
	c.ExecuteFun(engine)
}

func NewCustomCommand(execute func(engine *GameEngine)) *CustomCommand {
	return &CustomCommand{
		ExecuteFun: execute,
	}
}

type PlayerInputCommand struct {
	PlayerId PlayerId
	Cmd      protocol.ClientInputCommand
}

func (c PlayerInputCommand) Execute(engine *GameEngine) {
	playerInfo := engine.Players[c.PlayerId]

	if engine.Mod == SelectTeamMode {
		switch c.Cmd.Id {
		case protocol.InputCommandKind.ReadyToStart:
			playerInfo.IsReadyToStart = true
		case protocol.InputCommandKind.NotReadyToStart:
			playerInfo.IsReadyToStart = false
		}

		if playerInfo.IsReadyToStart {
			// If player is ready to start, he should be immovable
			return
		}
	}

	switch c.Cmd.Id {
	case protocol.InputCommandKind.MouseClick:
		x := c.Cmd.FloatArgs["x"]
		y := c.Cmd.FloatArgs["y"]
		fmt.Printf("processInputCommands: Click %f %f\n", x, y)
		AddBox(engine.World, x, y, 1, 1, 1, 0.5, 0.3)
	case protocol.InputCommandKind.MoveHero:
		playerBody := playerInfo.Body
		playerVel := playerBody.GetLinearVelocity()

		desiredVelX := playerVel.X
		desiredVelY := playerVel.Y

		moveKind := c.Cmd.IntArgs["kind"]
		direction := playerInfo.Direction
		switch moveKind {
		case protocol.MoveHeroKind.Right:
			desiredVelX = min(playerVel.X+settings.PlayerHorizontalAccelerationPerFrame, settings.PlayerMaxHorizontalSpeed)
			direction = protocol.DirectionKindRight
		case protocol.MoveHeroKind.Left:
			desiredVelX = max(playerVel.X-settings.PlayerHorizontalAccelerationPerFrame, -settings.PlayerMaxHorizontalSpeed)
			direction = protocol.DirectionKindLeft
		case protocol.MoveHeroKind.Up:
			if playerInfo.JumpCount > 0 {
				desiredVelY = settings.PlayerJumpSpeed
				playerInfo.JumpCount--
			}
		case protocol.MoveHeroKind.Down:
			playerInfo.MoveDownThrowPlatform = true
			desiredVelY = min(settings.PlayerDownSpeed, playerVel.Y)
		}

		velChangeX := desiredVelX - playerVel.X
		velChangeY := desiredVelY - playerVel.Y
		impulse := box2d.B2Vec2{
			X: playerBody.GetMass() * velChangeX,
			Y: playerBody.GetMass() * velChangeY,
		}
		playerBody.ApplyLinearImpulse(impulse, playerBody.GetWorldCenter(), true)
		playerInfo.Direction = direction
	case protocol.InputCommandKind.RotateHero:
		playerBody := playerInfo.Body
		kind := c.Cmd.IntArgs["kind"]
		var desiredSpeed float64
		switch kind {
		case protocol.RotateHeroKind.Left:
			desiredSpeed = settings.PlayerAngularSpeed
		case protocol.RotateHeroKind.Right:
			desiredSpeed = -settings.PlayerAngularSpeed
		}
		playerBody.SetAngularVelocity(desiredSpeed)
	case protocol.InputCommandKind.MakeShoot:
		weapon := playerInfo.Weapon
		weapon.Shoot(engine, playerInfo)
	case protocol.InputCommandKind.ChangeWeapon:
		weaponKind := protocol.WeaponKind(c.Cmd.IntArgs["kind"])
		weapon := CreateWeapon(weaponKind)
		playerInfo.Weapon = weapon
	}
	// TODO: remove
	engine.Players[c.PlayerId] = playerInfo
}

type CreatePlayerCommand struct {
	Nickname   string
	Team       protocol.TeamKind
	PlayerId   PlayerId
	PosX, PosY float64
}

func (c CreatePlayerCommand) Execute(engine *GameEngine) {
	// Hero body
	body := AddHero(engine.World, c.PosX, c.PosY, 0.8, 1, 1, 0.3, c.PlayerId)
	engine.Players[c.PlayerId] = &PlayerInfo{
		Nickname:              c.Nickname,
		Team:                  c.Team,
		Body:                  body,
		Direction:             protocol.DirectionKindRight,
		MoveDownThrowPlatform: false,
		JumpCount:             settings.PlayerMaxJumpCount,
		Weapon:                NewDefaultGun(),
	}
	fmt.Printf("createPlayerCommand: id=%d\n", c.PlayerId)
}

type CreateBulletCommand struct {
	PlayerId PlayerId
}

func (c CreateBulletCommand) Execute(engine *GameEngine) {
	playerBody := engine.Players[c.PlayerId].Body
	playerBody.GetPosition()
}

type RemoveBodyCommand struct {
	body *box2d.B2Body
}

func (c RemoveBodyCommand) Execute(engine *GameEngine) {
	engine.World.DestroyBody(c.body)
}

type ApplyImpulseCommand struct {
	body    *box2d.B2Body
	point   box2d.B2Vec2
	impulse box2d.B2Vec2
}

func (c ApplyImpulseCommand) Execute(engine *GameEngine) {
	c.body.ApplyLinearImpulse(c.impulse, c.point, true)
}
