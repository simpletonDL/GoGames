package engine

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/simpletonDL/GoGames/common/protocol"
)

type GameCommand interface {
	Execute(world *GameEngine)
}

type PlayerInputCommand struct {
	PlayerId PlayerId
	Cmd      protocol.ClientInputCommand
}

func (c PlayerInputCommand) Execute(engine *GameEngine) {
	switch c.Cmd.Id {
	case protocol.InputCommandKind.MouseClick:
		x := c.Cmd.FloatArgs["x"]
		y := c.Cmd.FloatArgs["y"]
		fmt.Printf("processInputCommands: Click %f %f\n", x, y)
		AddBox(engine.World, box2d.B2BodyType.B2_dynamicBody, x, y, 1, 1, 1, 1, 0.3)
	case protocol.InputCommandKind.MoveHero:
		playerBody := engine.Players[c.PlayerId].Body
		playerVel := playerBody.GetLinearVelocity()

		desiredVelX := playerVel.X
		desiredVelY := playerVel.Y

		moveKind := c.Cmd.IntArgs["kind"]
		switch moveKind {
		case protocol.MoveHeroKind.Right:
			desiredVelX = 5
		case protocol.MoveHeroKind.Left:
			desiredVelX = -5
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
	case protocol.InputCommandKind.MakeShoot:
		fmt.Printf("processInputCommands: Shoot\n")
		playerBody := engine.Players[c.PlayerId].Body
		playerPosition := playerBody.GetPosition()
		userData := playerBody.GetUserData().(BodyUserData)
		bullet := AddBullet(engine.World, playerPosition.X+userData.Width/1.5, playerPosition.Y, 0, 0.2, 0.2)
		bullet.SetLinearVelocity(box2d.B2Vec2{X: 15, Y: 0})
	}
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
