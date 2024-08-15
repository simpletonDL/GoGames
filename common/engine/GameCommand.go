package engine

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/GoGames/common/settings"
)

type GameCommand interface {
	Execute(engine *GameEngine)
}

type PlayerInputCommand struct {
	PlayerId PlayerId
	Cmd      protocol.ClientInputCommand
}

func (c PlayerInputCommand) Execute(engine *GameEngine) {
	playerInfo := engine.Players[c.PlayerId]
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
			direction = PlayerDirection.Right
		case protocol.MoveHeroKind.Left:
			desiredVelX = max(playerVel.X-settings.PlayerHorizontalAccelerationPerFrame, -settings.PlayerMaxHorizontalSpeed)
			direction = PlayerDirection.Left
		case protocol.MoveHeroKind.Up:
			desiredVelY = settings.PlayerJumpSpeed
			//case protocol.MoveHeroKind.Down:
			//	desiredVelY = -2
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
		fmt.Printf("processInputCommands: Shoot\n")
		playerBody := playerInfo.Body
		playerPosition := playerBody.GetPosition()
		bullet := AddBullet(engine.World, playerPosition.X, playerPosition.Y, 0, 0.2, 0.2, playerBody)
		bulletRotation := box2d.MakeB2RotFromAngle(playerBody.GetAngle())
		bulletVec := box2d.MakeB2Vec2(bulletRotation.C, bulletRotation.S)
		bulletVec.OperatorScalarMulInplace(15.0)
		if playerInfo.Direction == PlayerDirection.Left {
			bulletVec.OperatorScalarMulInplace(-1.0)
		}
		bullet.SetLinearVelocity(bulletVec)
	}
	engine.Players[c.PlayerId] = playerInfo
}

type CreatePlayerCommand struct {
	PlayerId   PlayerId
	PosX, PosY float64
}

func (c CreatePlayerCommand) Execute(engine *GameEngine) {
	// Hero body
	body := AddHero(engine.World, 2, 15, 0.8, 1, 1, 0.3)
	engine.Players[c.PlayerId] = PlayerInfo{Body: body}
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
