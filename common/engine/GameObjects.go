package engine

import (
	"github.com/simpletonDL/GoGames/common/protocol"
	"github.com/simpletonDL/box2d"
)

type BodyUserData interface {
	GetWidth() float64
	GetHeight() float64
	GetKind() protocol.BodyKind
}

type DefaultBodyUserData struct {
	Width  float64
	Height float64
	Kind   protocol.BodyKind
}

func (d DefaultBodyUserData) GetWidth() float64 {
	return d.Width
}

func (d DefaultBodyUserData) GetHeight() float64 {
	return d.Height
}

func (d DefaultBodyUserData) GetKind() protocol.BodyKind {
	return d.Kind
}

type PlayerUserData struct {
	DefaultBodyUserData
	HeroId PlayerId
}

type BulletUserData struct {
	DefaultBodyUserData
	Owner       *box2d.B2Body
	ImpactForce float64
}

func SetBodyUserData(body *box2d.B2Body, data BodyUserData) {
	body.SetUserData(data)
}

func GetBodyUserData(body *box2d.B2Body) BodyUserData {
	return body.GetUserData().(BodyUserData)
}

//type BodyUserData struct {
//	Width  float64
//	Height float64
//	Kind   protocol.BodyKind
//	// For case when Kind is Hero
//	HeroId PlayerId
//	// Bodies don't collide with their owner (for skip collisions between player and created by him bullets)
//	Owner *box2d.B2Body
//}

func AddBox(world *box2d.B2World, x float64, y float64, angel float64, width float64, height float64, density float64, friction float64) *box2d.B2Body {
	body := addRectangle(world, box2d.B2BodyType.B2_dynamicBody, x, y, angel, width, height, density, friction)
	SetBodyUserData(body, DefaultBodyUserData{Width: width, Height: height, Kind: protocol.BodyKindBox})
	return body
}

func AddWeaponBox(world *box2d.B2World, x float64, y float64, angel float64, width float64, height float64, density float64, friction float64) *box2d.B2Body {
	body := addRectangle(world, box2d.B2BodyType.B2_dynamicBody, x, y, angel, width, height, density, friction)
	SetBodyUserData(body, DefaultBodyUserData{Width: width, Height: height, Kind: protocol.BodyKindWeaponBox})
	return body
}

func AddPlatform(world *box2d.B2World, x float64, y float64, angel float64, width float64, height float64, density float64, friction float64) *box2d.B2Body {
	body := addRectangle(world, box2d.B2BodyType.B2_staticBody, x, y, angel, width, height, density, friction)
	SetBodyUserData(body, DefaultBodyUserData{Width: width, Height: height, Kind: protocol.BodyKindPlatform})
	return body
}

func AddHero(world *box2d.B2World, x float64, y float64, width float64, height float64, density float64, friction float64, id PlayerId) *box2d.B2Body {
	hero := addRectangle(world, box2d.B2BodyType.B2_dynamicBody, x, y, 0, width, height, density, friction)
	SetBodyUserData(hero, PlayerUserData{
		DefaultBodyUserData: DefaultBodyUserData{Width: width, Height: height, Kind: protocol.BodyKindHero},
		HeroId:              id,
	})

	hero.SetFixedRotation(true)
	return hero
}

func AddBullet(
	world *box2d.B2World, x float64, y float64, angel float64, width float64, height float64,
	owner *box2d.B2Body, impactForce float64,
) *box2d.B2Body {
	bullet := addRectangle(world, box2d.B2BodyType.B2_kinematicBody, x, y, 0, width, height, 1, 1)
	SetBodyUserData(bullet, BulletUserData{
		DefaultBodyUserData: DefaultBodyUserData{Width: width, Height: height, Kind: protocol.BodyKindBullet},
		Owner:               owner,
		ImpactForce:         impactForce,
	})
	bullet.SetBullet(true)
	return bullet
}

func addRectangle(world *box2d.B2World, bodyType uint8, x float64, y float64, angel float64, width float64, height float64, density float64, friction float64) *box2d.B2Body {
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
	return body
}

func NewWorld(gravityX float64, gravityY float64) *box2d.B2World {
	gravity := box2d.MakeB2Vec2(gravityX, gravityY)
	world := box2d.MakeB2World(gravity)
	return &world
}

func createInitialWorld() *box2d.B2World {
	// Create default map

	world := NewWorld(0, -20)
	// Platforms
	AddPlatform(world, 12, 1, 0, 22, 1, 0, 1)
	AddPlatform(world, 12, 4.5, 0, 18, 1, 0, 1)
	AddPlatform(world, 4, 8, 0, 7, 1, 0, 1)
	AddPlatform(world, 20, 8, 0, 7, 1, 0, 1)
	AddPlatform(world, 12, 12, 0, 7, 1, 0, 1)

	// Dynamic body
	AddBox(world, 8, 15, 1, 1, 1, 1, 0.3)

	return world
}
