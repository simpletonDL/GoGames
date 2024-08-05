package engine

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/simpletonDL/GoGames/common/protocol"
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
		AddBox(engine.World, box2d.B2BodyType.B2_dynamicBody, x, y, 1, 1, 1, 1, 0.3)
	case protocol.InputCommandKind.MoveHero:
		playerBody := playerInfo.Body
		playerVel := playerBody.GetLinearVelocity()

		desiredVelX := playerVel.X
		desiredVelY := playerVel.Y

		moveKind := c.Cmd.IntArgs["kind"]
		direction := playerInfo.Direction
		switch moveKind {
		case protocol.MoveHeroKind.Right:
			desiredVelX = 5
			direction = PlayerDirection.Right
		case protocol.MoveHeroKind.Left:
			desiredVelX = -5
			direction = PlayerDirection.Left
		case protocol.MoveHeroKind.Up:
			desiredVelY = 10
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
	case protocol.InputCommandKind.MakeShoot:
		fmt.Printf("processInputCommands: Shoot\n")
		playerBody := playerInfo.Body
		playerPosition := playerBody.GetPosition()
		bullet := AddBullet(engine.World, playerPosition.X, playerPosition.Y, 0, 0.2, 0.2, playerBody)
		velX := 15.0
		if playerInfo.Direction == PlayerDirection.Left {
			velX *= -1
		}
		bullet.SetLinearVelocity(box2d.B2Vec2{X: velX, Y: 0})
	}
	engine.Players[c.PlayerId] = playerInfo
}

type CreatePlayerCommand struct {
	PlayerId   PlayerId
	PosX, PosY float64
}

func (c CreatePlayerCommand) Execute(engine *GameEngine) {
	// Hero body
	body := AddHero(engine.World, 2, 15, 0.6, 1, 1, 0.3)
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
