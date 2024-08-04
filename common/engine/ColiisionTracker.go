package engine

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/simpletonDL/GoGames/common/protocol"
)

type CollisionTracker struct {
	engine *GameEngine
}

func NewCollisionTracker(engine *GameEngine) CollisionTracker {
	return CollisionTracker{engine: engine}
}

func (c CollisionTracker) BeginContact(contact box2d.B2ContactInterface) {}

func (c CollisionTracker) EndContact(contact box2d.B2ContactInterface) {}

func (c CollisionTracker) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	bodyA := contact.GetFixtureA().GetBody()
	bodyB := contact.GetFixtureB().GetBody()
	userDataA := bodyA.GetUserData().(BodyUserData)
	userDataB := bodyB.GetUserData().(BodyUserData)

	// At lease one body should be a bullet
	if userDataA.Kind != protocol.BodyKind.Bullet && userDataB.Kind != protocol.BodyKind.Bullet {
		return
	}

	// Make sure that bodyA is a bullet
	if userDataB.Kind == protocol.BodyKind.Bullet {
		bodyA, bodyB = bodyB, bodyA
		userDataA, userDataB = userDataB, userDataA
	}

	if bodyB == userDataA.Owner {
		// bullets should not contact with their owners
		contact.SetEnabled(false)
		return
	}

	fmt.Printf("Bullet(%f, %f) contact\n", bodyA.GetPosition().X, bodyA.GetPosition().Y)
	contact.SetEnabled(false)
	c.engine.ScheduleCommand(RemoveBodyCommand{body: bodyA})
	if userDataB.Kind == protocol.BodyKind.Bullet {
		c.engine.ScheduleCommand(RemoveBodyCommand{body: bodyB})
	} else {
		var worldManifold box2d.B2WorldManifold
		contact.GetWorldManifold(&worldManifold)
		if contact.GetManifold().PointCount > 0 {
			collisionPoint := worldManifold.Points[0]
			fmt.Printf("Collision point: (%f, %f)\n", collisionPoint.X, collisionPoint.Y)
			fmt.Printf("Body world center: (%f, %f)", bodyB.GetWorldCenter().X, bodyB.GetWorldCenter().Y)
			c.engine.ScheduleCommand(ApplyImpulseCommand{
				body:    bodyB,
				point:   collisionPoint,
				impulse: box2d.B2Vec2{X: 7, Y: 0},
			})
		}
	}
}

func (c CollisionTracker) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
}
