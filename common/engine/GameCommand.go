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
