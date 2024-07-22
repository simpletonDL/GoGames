package engine

import (
	"github.com/ByteArena/box2d"
	"github.com/simpletonDL/GoGames/common/protocol"
)

type BodyUserData struct {
	Width  float64
	Height float64
}

func B2WorldToGameState(world *box2d.B2World) protocol.GameState {
	gameObjects := make([]protocol.GameObject, world.GetBodyCount())
	for body := world.GetBodyList(); body != nil; body = body.M_next {
		data := body.GetUserData().(BodyUserData)
		object := protocol.GameObject{
			XPos:      body.GetPosition().X,
			YPos:      body.GetPosition().Y,
			Angel:     body.GetAngle(),
			ImageKind: protocol.BoxKind,
			Width:     data.Width,
			Height:    data.Height,
		}
		gameObjects = append(gameObjects, object)
	}
	return protocol.GameState{
		Objects: gameObjects,
	}
}

func NewWorld(gravityX float64, gravityY float64) *box2d.B2World {
	gravity := box2d.MakeB2Vec2(gravityX, gravityY)
	world := box2d.MakeB2World(gravity)
	return &world
}

func AddBox(world *box2d.B2World, bodyType uint8, x float64, y float64, angel float64, width float64, height float64, density float64, friction float64) *box2d.B2Body {
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = bodyType
	bodyDef.Position.Set(x, y)
	bodyDef.Angle = angel
	body := world.CreateBody(&bodyDef)

	boxShape := box2d.MakeB2PolygonShape()
	boxShape.SetAsBox(width/2, height/2)
	fixtureDef := box2d.MakeB2FixtureDef()
	fixtureDef.Shape = &boxShape
	fixtureDef.Density = density
	fixtureDef.Friction = friction
	body.CreateFixtureFromDef(&fixtureDef)

	body.SetUserData(BodyUserData{Width: width, Height: height})
	return body
}
