package engine

import (
	"github.com/ByteArena/box2d"
	"github.com/simpletonDL/GoGames/common/protocol"
)

type BodyUserData struct {
	Width  float64
	Height float64
	Kind   uint8
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

	body.SetUserData(BodyUserData{Width: width, Height: height, Kind: protocol.BodyKind.Box})
	return body
}

func AddHero(world *box2d.B2World, x float64, y float64, width float64, height float64, density float64, friction float64) *box2d.B2Body {
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.Position.Set(x, y)
	body := world.CreateBody(&bodyDef)

	boxShape := box2d.MakeB2PolygonShape()
	boxShape.SetAsBox(width/2, height/2)
	fixtureDef := box2d.MakeB2FixtureDef()
	fixtureDef.Shape = &boxShape
	fixtureDef.Density = density
	fixtureDef.Friction = friction
	body.CreateFixtureFromDef(&fixtureDef)

	body.SetUserData(BodyUserData{Width: width, Height: height, Kind: protocol.BodyKind.Hero})
	return body
}

func NewWorld(gravityX float64, gravityY float64) *box2d.B2World {
	gravity := box2d.MakeB2Vec2(gravityX, gravityY)
	world := box2d.MakeB2World(gravity)
	return &world
}

func createInitialWorld() *box2d.B2World {
	world := NewWorld(0, -20)
	// Ground body
	AddBox(world, box2d.B2BodyType.B2_staticBody, 8, 1, 0, 16, 2, 0, 1)

	// Dynamic body
	AddBox(world, box2d.B2BodyType.B2_dynamicBody, 8, 15, 1, 1, 1, 1, 0.3)

	return world
}
